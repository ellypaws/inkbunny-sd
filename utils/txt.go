package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ellypaws/inkbunny-sd/entities"
)

type Params map[string]PNGChunk
type PNGChunk map[string]string

type Config struct {
	Text          string
	KeyCondition  func(string) bool
	SkipCondition func(string) bool
	Filename      string
}

type Processor func(...func(*Config)) (Params, error)

const (
	Parameters     = "parameters"
	Postprocessing = "postprocessing"
	Extras         = "extras"
	Objects        = "objects"
	Caption        = "caption"

	IDAutoSnep    = 1004248
	IDDruge       = 151203
	IDArtieDragon = 1190392
	IDAIBean      = 147301
	IDFairyGarden = 215070
	IDCirn0       = 177167
	IDHornybunny  = 12499
	IDNeoncortex  = 14603
	IDMethuzalach = 1089071
	IDRNSDAI      = 1188211
	IDSoph        = 229969
	IDNastAI      = 1101622
)

// AutoSnep is a Processor that parses yaml like raw txt where each two spaces is a new dict
// It's mostly seen in multi-chunk parameter output from AutoSnep
func AutoSnep(opts ...func(*Config)) (Params, error) {
	var c Config
	for _, f := range opts {
		f(&c)
	}
	var chunks Params = make(Params)
	scanner := bufio.NewScanner(strings.NewReader(c.Text))

	const Software = "Software"

	var png string
	var key string
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		line := scanner.Text()

		indentLevel := len(line) - len(strings.TrimLeft(line, " "))

		switch indentLevel {
		case 0:
			if strings.HasSuffix(line, ":") {
				png = c.Filename + strings.TrimSuffix(line, ":")
				chunks[png] = make(PNGChunk)
			}
		case 2: // PNG text chunks:
			if strings.TrimSpace(line) == "PNG text chunks:" {
				chunks[png] = make(PNGChunk)
			}
		case 4: // parameters:
			key = strings.TrimSpace(strings.TrimSuffix(line, ":"))
		case 6:
			if len(chunks[png][key]) > 0 {
				chunks[png][key] += "\n"
			}
			chunks[png][key] += line[6:]
		default:

		}
	}

	if len(chunks) == 0 {
		return nil, errors.New("no chunks found")
	}

	return chunks, nil
}

var seedLine = regexp.MustCompile(`seed: (\d+)`)

func Cirn0(opts ...func(*Config)) (Params, error) {
	var c Config
	for _, f := range opts {
		f(&c)
	}
	var chunks Params = make(Params)
	scanner := bufio.NewScanner(strings.NewReader(c.Text))

	var steps, sampler, cfg, model string
	var key string
	var lastKey string
	var foundSeed bool
	var foundNegative bool
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)

		switch {
		case len(line) == 0:
			continue

		case strings.HasPrefix(line, "sampler:"):
			sampler = strings.TrimPrefix(line, "sampler: ")
		case strings.HasPrefix(line, "cfg:"):
			cfg = strings.TrimPrefix(line, "cfg: ")
		case strings.HasPrefix(line, "steps:"):
			steps = strings.TrimPrefix(line, "steps: ")
		case strings.HasPrefix(line, "model:"):
			model = strings.TrimPrefix(line, "model: ")

		case strings.HasPrefix(line, "==="):
			key = c.Filename + line
			chunks[key] = make(PNGChunk)

			lastKey = strings.Trim(line, "= ")

		case strings.HasPrefix(line, "---"):
			if len(lastKey) == 0 {
				lastKey = "unknown"
			}
			key = fmt.Sprintf("--- %s%s ---", lastKey, strings.Trim(line, "- "))
			chunks[key] = make(PNGChunk)

		case len(key) == 0:
			continue

		case strings.HasPrefix(line, "keyword prompt:"):
			continue

		case foundNegative:
			foundNegative = false
			chunks[key][Parameters] += fmt.Sprintf("\nNegative Prompt: %s", line)

		case strings.HasPrefix(line, "negative prompt:"):
			foundNegative = true
			continue

		case seedLine.MatchString(line):
			chunks[key][Parameters] += fmt.Sprintf(
				"\nSteps: %s, Sampler: %s, CFG scale: %s, Seed: %s, Model: %s",
				steps,
				sampler,
				cfg,
				seedLine.FindStringSubmatch(line)[1],
				model,
			)
			key = ""

		case foundSeed:
			foundSeed = false
			chunks[key][Parameters] += fmt.Sprintf(
				"\nSteps: %s, Sampler: %s, CFG scale: %s, Seed: %s, Model: %s",
				steps,
				sampler,
				cfg,
				line,
				model,
			)
			key = ""

		case strings.HasPrefix(line, "seed:"):
			foundSeed = true

		default:
			if len(chunks[key][Parameters]) > 0 {
				chunks[key][Parameters] += "\n"
			}
			chunks[key][Parameters] += line
		}
	}

	if len(chunks) == 0 {
		return nil, errors.New("no chunks found")
	}

	return chunks, nil
}

// SophStartKey is intended to check if the content contains a key in the format of ./file: for IDSoph.
// Otherwise, use Common and UseSoph to extract PNGChunk as normal.
var SophStartKey = regexp.MustCompile(`(?m)^\./[^:]*:$`)

// SophStartInvokeAI checks if the content contains any JSON objects for IDSoph.
// Otherwise, use Common and UseSoph to extract PNGChunk as normal.
var SophStartInvokeAI = regexp.MustCompile(`(?m)^{$`)

// Soph is a Processor that parses the InvokeAI format by IDSoph.
// Check if the content follows the InvokeAI format using SophStartKey.
func Soph(opts ...func(*Config)) (map[string]entities.TextToImageRequest, error) {
	var c Config
	for _, f := range opts {
		f(&c)
	}
	if len(c.Text) == 0 {
		return nil, errors.New("empty text")
	}

	scanner := bufio.NewScanner(bytes.NewReader([]byte(c.Text)))

	var invokeAI = make(map[string]entities.InvokeAI)
	var invokeBuffer bytes.Buffer
	var counter int
	var key string
	for scanner.Scan() {
		line := scanner.Bytes()
		switch {
		case len(line) == 0:
			continue
		case bytes.HasPrefix(scanner.Bytes(), []byte("./")):
			key = string(line)
			continue
		case bytes.Equal(line, []byte(`{`)):
			invokeBuffer.Write(line)
		case bytes.Equal(line, []byte(`}`)):
			invokeBuffer.Write(line)
			var i entities.InvokeAI
			if err := json.Unmarshal(invokeBuffer.Bytes(), &i); err != nil {
				return nil, err
			}
			if len(key) == 0 {
				counter++
				key = fmt.Sprintf("%s_%d", c.Filename, counter)
			}
			invokeAI[key] = i
			invokeBuffer.Reset()
			key = ""
		case invokeBuffer.Len() > 0:
			invokeBuffer.Write(line)
		default:
			continue
		}
	}

	if len(invokeAI) == 0 {
		return nil, errors.New("no chunks found")
	}

	var objects = make(map[string]entities.TextToImageRequest)
	for k, v := range invokeAI {
		objects[k] = v.Convert()
	}

	return objects, nil
}

func Sequential(opts ...func(*Config)) (Params, error) {
	var c Config
	for _, f := range opts {
		f(&c)
	}
	if len(c.Text) == 0 {
		return nil, errors.New("empty text")
	}

	var chunks Params = make(Params)
	var page int
	var key = fmt.Sprintf("%s-%d", c.Filename, page)

	scanner := bufio.NewScanner(strings.NewReader(c.Text))
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		line := scanner.Text()

		switch {
		case len(line) == 0:
			continue
		case c.SkipCondition != nil && c.SkipCondition(line):
			continue
		case strings.HasPrefix(line, "Steps:"):
			chunks[key][Parameters] += "\n" + line
			page++
			key = fmt.Sprintf("%s-%d", c.Filename, page)
		default:
			if chunks[key] == nil {
				chunks[key] = make(PNGChunk)
			}
			if len(chunks[key][Parameters]) > 0 {
				chunks[key][Parameters] += "\n"
			}
			chunks[key][Parameters] += line
		}
	}

	if len(chunks) == 0 {
		return nil, errors.New("no chunks found")
	}

	return chunks, nil
}

var drugeMatchDigit = regexp.MustCompile(`(?m)^\d+`)

func UseDruge() func(*Config) {
	return func(c *Config) {
		c.KeyCondition = func(line string) bool {
			return drugeMatchDigit.MatchString(line)
		}
		c.Filename = "druge_"
		if !drugeMatchDigit.MatchString(c.Text) {
			c.Text = "1\n" + c.Text
		}
	}
}

func UseArtie() func(*Config) {
	return func(c *Config) {
		c.KeyCondition = func(line string) bool {
			return strings.HasSuffix(line, "Image")
		}
		c.Filename = "artiedragon_"
	}
}

var aiBeanKey = regexp.MustCompile(`(?i)^(image )?\d+`)

func UseAIBean() func(*Config) {
	return func(c *Config) {
		c.KeyCondition = func(line string) bool {
			return aiBeanKey.MatchString(line)
		}
		c.Filename = "AIBean_"
		c.SkipCondition = func(line string) bool {
			return line == "parameters"
		}
		if aiBeanKey.MatchString(c.Text) {
			return
		}
		if strings.HasPrefix(c.Text, "parameters") {
			c.Text = strings.Replace(c.Text, "parameters", "1", 1)
			return
		}
		c.Text = "1\n" + c.Text
	}
}

func UseFairyGarden() func(*Config) {
	return func(c *Config) {
		c.KeyCondition = func(line string) bool {
			return strings.HasPrefix(line, "photo")
		}
		c.Filename = "fairygarden_"
		// prepend "photo 1" to the input in case it's missing
		c.Text = "photo 1\n" + c.Text
	}
}

func UseCirn0() func(*Config) {
	return func(c *Config) {
		c.KeyCondition = func(line string) bool {
			return strings.HasPrefix(line, "===")
		}
		c.Filename = "cirn0_"

		var part string
		lines := strings.Split(c.Text, "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "=== #") {
				part = strings.TrimPrefix(line, "=== #")
			}
			if strings.HasPrefix(line, "---") {
				lines[i] = fmt.Sprintf("=== Part #%s", part)
			}
		}
		c.Text = strings.Join(lines, "\n")
	}
}

func UseHornybunny() func(*Config) {
	return func(c *Config) {
		c.Text = "(1)\n" + c.Text
		c.KeyCondition = func(line string) bool {
			return regexp.MustCompile(`^\(\d+\)$`).MatchString(line)
		}
		c.Filename = "Hornybunny_"
		// c.Text = strings.ReplaceAll(c.Text, "----", "")
		// c.Text = strings.ReplaceAll(c.Text, "Original generation details", "")
		// c.Text = strings.ReplaceAll(c.Text, "Upscaling details", "")
		c.Text = strings.ReplaceAll(c.Text, "Positive Prompt: ", "")
		c.Text = strings.ReplaceAll(c.Text, "Other details: ", "")
		c.SkipCondition = func(line string) bool {
			switch line {
			case "----":
				return true
			case "==========":
				return true
			case "Original generation details":
				return true
			case "Upscaling details":
				return true
			default:
				return false
			}
		}
	}
}

var (
	methuzalachModel    = regexp.MustCompile(`Model: [^\n]+`)
	methuzalachNegative = regexp.MustCompile(`Negative prompts:\s*`)
	methuzalachSeed     = regexp.MustCompile(`\s*Seed: \D*?[,\s]`)
	methuzalachSteps    = regexp.MustCompile(`.*(Steps: \d+[^\n]*)`)
)

func UseMethuzalach() func(*Config) {
	return func(c *Config) {
		c.KeyCondition = func(line string) bool {
			return strings.HasPrefix(line, "Image")
		}

		model := methuzalachModel.FindString(c.Text)
		c.Text = methuzalachNegative.ReplaceAllString(c.Text, "Negative Prompt: ")
		c.Text = methuzalachSeed.ReplaceAllString(c.Text, "")
		c.Text = methuzalachSteps.ReplaceAllString(c.Text, `$1 `+model)
	}
}

func UseSoph() func(*Config) {
	return func(c *Config) {
		c.KeyCondition = func(line string) bool {
			return strings.HasPrefix(line, c.Filename)
		}
		c.Text = c.Filename + "\n" + c.Text
		c.SkipCondition = func(line string) bool {
			return strings.HasPrefix(line, "<comment: ")
		}
	}
}

func WithString(s string) func(*Config) {
	return func(c *Config) {
		c.Text = s
	}
}

func WithBytes(b []byte) func(*Config) {
	return func(c *Config) {
		c.Text = string(b)
	}
}

func WithConfig(config Config) func(*Config) {
	return func(c *Config) {
		*c = config
	}
}

func WithFilename(filename string) func(*Config) {
	return func(c *Config) {
		c.Filename = filename
	}
}

func WithKeyCondition(f func(string) bool) func(*Config) {
	return func(c *Config) {
		c.KeyCondition = f
	}
}

func Common(opts ...func(*Config)) (Params, error) {
	var c Config
	for _, f := range opts {
		f(&c)
	}
	if c.KeyCondition == nil {
		return nil, errors.New("condition for key is not set")
	}
	var chunks Params = make(Params)
	scanner := bufio.NewScanner(strings.NewReader(c.Text))

	var (
		key            string
		negativePrompt string
		extra          string

		objectStart string
		objectKey   string
		objectCount int
	)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		line := scanner.Text()

		if c.SkipCondition != nil && c.SkipCondition(line) {
			continue
		}
		if len(line) == 0 {
			continue
		}
		if c.KeyCondition(line) {
			key = c.Filename + line
			chunks[key] = make(PNGChunk)
			continue
		}
		if len(key) == 0 {
			continue
		}
		if len(negativePrompt) > 0 {
			if !negativeHasText.MatchString(negativePrompt) {
				chunks[key][Parameters] += line
				continue
			}
			chunks[key][Parameters] += "\n" + line
			if stepsStart.MatchString(line) {
				negativePrompt = ""
				key = ""
			}
			continue
		}
		if negativeStart.MatchString(line) {
			negativePrompt = line
			chunks[key][Parameters] += "\n" + line
			continue
		}
		if len(objectStart) > 0 {
			chunks[key][objectKey] += "\n" + line
			if objectStart == "{" && line == "}" {
				objectCount++
				objectKey = ""
				objectStart = ""
			}
			if objectStart == "[" && line == "]" {
				objectCount++
				objectKey = ""
				objectStart = ""
			}
		}
		if len(extra) > 0 {
			chunks[key][extra] += line
			extra = ""
			continue
		}
		switch line {
		case "{", "[":
			objectStart = line
			objectKey = fmt.Sprintf("%s_%d", Objects, objectCount)
			chunks[key][objectKey] += line
		case Postprocessing:
			extra = Postprocessing
			continue
		case Extras:
			extra = Extras
			continue
		}
		if len(chunks[key][Parameters]) > 0 {
			chunks[key][Parameters] += "\n"
		}
		chunks[key][Parameters] += line
	}

	if len(chunks) == 0 {
		return nil, errors.New("no chunks found")
	}

	return chunks, nil
}

func ParseParams(p Params) map[string]entities.TextToImageRequest {
	var request map[string]entities.TextToImageRequest
	for file, chunk := range p {
		if params, ok := chunk[Parameters]; ok {
			r, err := ParameterHeuristics(params)
			if err != nil {
				continue
			}
			if request == nil {
				request = make(map[string]entities.TextToImageRequest)
			}
			request[file] = r
		}
	}
	return request
}
