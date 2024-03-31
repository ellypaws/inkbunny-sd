package utils

import (
	"bufio"
	"errors"
	"github.com/ellypaws/inkbunny-sd/entities"
	"io"
	"os"
	"strconv"
	"strings"
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

const Parameters = "parameters"

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
		line := scanner.Text()

		indentLevel := len(line) - len(strings.TrimLeft(line, " "))

		switch indentLevel {
		case 0:
			if strings.HasSuffix(line, ":") {
				png = "AutoSnep_" + strings.TrimSuffix(line, ":")
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

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return chunks, nil
}

func UseDruge() func(*Config) {
	return func(c *Config) {
		c.KeyCondition = func(line string) bool {
			_, err := strconv.Atoi(line)
			return err == nil
		}
		c.Filename = "druge_"
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

func UseAIBean() func(*Config) {
	return func(c *Config) {
		c.KeyCondition = func(line string) bool {
			_, err := strconv.Atoi(line)
			return err == nil
		}
		c.Filename = "AIBean_"
		c.SkipCondition = func(line string) bool {
			return line == "parameters"
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

	var key string
	var foundNegative bool
	for scanner.Scan() {
		line := scanner.Text()
		if c.SkipCondition != nil && c.SkipCondition(line) {
			continue
		}
		if foundNegative {
			chunks[key][Parameters] += "\n" + line
			foundNegative = false
			key = ""
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
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "Negative Prompt:") {
			foundNegative = true
			chunks[key][Parameters] += "\n" + line
			continue
		}
		if len(chunks[key][Parameters]) > 0 {
			chunks[key][Parameters] += "\n"
		}
		chunks[key][Parameters] += line
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

func FileToRequests(file string, processor Processor, opts ...func(*Config)) (map[string]entities.TextToImageRequest, error) {
	p, err := fileToParams(file, processor, opts...)
	if err != nil {
		return nil, err
	}
	return ParseParams(p), nil
}

// fileToParams reads the file and returns the params using a Processor
func fileToParams(file string, processor Processor, opts ...func(*Config)) (Params, error) {
	f, err := fileToBytes(file)
	if err != nil {
		return nil, err
	}
	opts = append(opts, WithFilename(file))
	opts = append(opts, WithBytes(f))
	return processor(opts...)
}

// fileToBytes reads the file and returns the content as a byte slice
func fileToBytes(file string) ([]byte, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(f)
}
