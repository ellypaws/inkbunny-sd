package utils

import (
	"bytes"
	"regexp"
)

var (
	// Patterns are the regexp.Regexp patterns for DescriptionHeuristics to get the parameters from the description.
	// The keys are the names of the capture groups in the pattern. Currently being used in ExtractAll
	Patterns = map[string]*regexp.Regexp{
		"steps":      regexp.MustCompile(`(?i)steps[:\s-]+(?P<steps>\d+)`),
		"sampler":    regexp.MustCompile(`(?i)sampl(?:er|ing method)[:\s-]+(?P<sampler>[\w+ ]+)`),
		"cfg":        regexp.MustCompile(`(?i)(?:cfg(?: scale)?|scale)[:\s-]+(?P<cfg>[\d.]+)`),
		"seed":       regexp.MustCompile(`(?i)seeds?[:\s-]+?(?P<seed>-?\d+)`),
		"width":      regexp.MustCompile(`(?i)size[:\s-]+(?P<width>\d+)x\d+`),
		"height":     regexp.MustCompile(`(?i)size[:\s-]+\d+x(?P<height>\d+)`),
		"hash":       regexp.MustCompile(`(?i)model hash[:\s-]+(?P<hash>\w+)`),
		"model":      regexp.MustCompile(`(?i)(?:model|checkpoint)s?[:\s-]+(?P<model>[^,\n]+)`),
		"denoising":  regexp.MustCompile(`(?i)denoising strength[:\s-]+(?P<denoising>[\d.]+)`),
		"loraHashes": regexp.MustCompile(loraHashes),
		"tiHashes":   regexp.MustCompile(tiHashes),
		"version":    regexp.MustCompile(`(?i)version[:\s-]+(?P<version>v[\w.-]+)`),
	}

	// RNSDAIPatterns are preset regexp.Regexp patterns for IDRNSDAI
	RNSDAIPatterns = map[string]*regexp.Regexp{
		"model":    regexp.MustCompile(`(?i)model[\s•]*\[b](?P<model>[^[]+)\[/b]`),
		"seed":     regexp.MustCompile(`(?i)seeds[\s•]*\[(?P<seed>\d+)]?`),
		"version":  regexp.MustCompile(`(?i)image generator[\s•]*[^:]+: v[\d.]+`),
		"prompt":   regexp.MustCompile(`(?is)positive prompt:\[/b]\n(?P<prompt>.*?)\n\[/q]`),
		"negative": regexp.MustCompile(`(?is)negative prompt:\[/b]\n(?P<negative>.*?)\n\[/q]`),
	}

	allParams = regexp.MustCompile(`\s*(\w[\w \-/]+):\s*("(?:\\.|[^\\"])+"|[^,]*)(?:,|$)`)

	positivePattern = regexp.MustCompile(`(?ims)^(?:(?:primary |pos(?:itive)? )?prompts?:?)\s*(.+?)\s*negative`)
	positiveEnd     = regexp.MustCompile(`(?ims)^(?:(?:primary |pos(?:itive)? )?prompts?:?)\s*(.+)`)
	negativePattern = regexp.MustCompile(`(?ims)^(?:(?:neg(?:ative)?)(?: prompts?)?:?)\s*(.+?)\s*(?:steps|sampler|model|seed|cfg)`)
	negativeEnd     = regexp.MustCompile(`(?ims)^(?:(?:neg(?:ative)?)(?: prompts?)?:?)\s*(.+)`)
	negativeStart   = regexp.MustCompile(`(?i)^negative(?: prompts?)?:\s*`)
	ParametersStart = regexp.MustCompile(`(?ims)^(parameters\n.*)`)

	bbCode = regexp.MustCompile(`\[/?[^]]+]`)
	emojis = regexp.MustCompile(`[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{1F700}-\x{1F77F}\x{1F780}-\x{1F7FF}\x{1F800}-\x{1F8FF}\x{1F900}-\x{1F9FF}\x{1FA00}-\x{1FA6F}\x{1FA70}-\x{1FAFF}\x{2700}-\x{27BF}\x{2600}-\x{26FF}\x{1F1E0}-\x{1F1FF}]`)

	negativeHasText = regexp.MustCompile(`(?i)^negative prompt: ?\S`)
	stepsStart      = regexp.MustCompile(`(?i)^steps: ?\d`)
	StepsStart      = regexp.MustCompile(`(?im)^Steps: ?\d+, Sampler:`)

	extractJson    = regexp.MustCompile(`(?ms){.*}`)
	removeComments = regexp.MustCompile(`(?m)//.*$`)
	fixParentheses = regexp.MustCompile(`\\+([()])`)

	newLineFix       = []byte("\n")
	quoteFix         = []byte(`\"`)
	stringExtraction = regexp.MustCompile(`(?s)(:\s+")(.*?)(",\n)|(?s)(:\s+")(.*?)(".*?\n)`)

	parenthesesReplacement = []byte(`\\$1`)
	newLineReplacement     = []byte(`\n`)
	quoteReplacement       = []byte(`"`)
)

const (
	loraHashes = `(?i)lora hashes:? "(?P<lora>[^"]+)"`
	loraEntry  = `(?i)(?P<key>\w+): (?P<value>\w+)`
	tiHashes   = `(?i)ti hashes:? "(?P<ti>[^"]+)"`
	tiEntry    = `(?i)(?P<key>\w+): (?P<value>\w+)`
)

func RemoveBBCode(s string) string {
	return bbCode.ReplaceAllString(s, "")
}

func RemoveEmojis(s string) string {
	return emojis.ReplaceAllString(s, "")
}

func CleanText(s string) string {
	s = RemoveBBCode(s)
	s = RemoveEmojis(s)
	return s
}

type ExtractResult map[string]string

func ExtractAll(s string, reg map[string]*regexp.Regexp) ExtractResult {
	var result = make(ExtractResult)

	for key, r := range reg {
		result[key] = Extract(s, r)
	}

	return result
}

// ExtractKeys extracts keys and values from a string using regex.
// Usually only the last line is passed here
//
//	`\s*(\w[\w \-/]+):\s*("(?:\\.|[^\\"])+"|[^,]*)(?:,|$)`
//
// [Key: value], [Key: value], [Key: "key: value, key:value"]
func ExtractKeys(line string) ExtractResult {
	var result = make(ExtractResult)

	matches := allParams.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		result[match[1]] = match[2]
	}

	return result
}

// ExtractDefaultKeys can use a default mapping to store the results to.
// Usually only the last line is passed here
func ExtractDefaultKeys(line string, defaultResults ExtractResult) ExtractResult {
	if defaultResults == nil {
		defaultResults = make(ExtractResult)
	}

	matches := allParams.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		defaultResults[match[1]] = match[2]
	}

	return defaultResults
}

func Extract(s string, r *regexp.Regexp) string {
	match := r.FindStringSubmatch(s)

	for i, name := range match {
		if i != 0 {
			return name
		}
	}

	return ""
}

func ExtractJson(content []byte) []byte {
	// Extracts everything from `{.*}`
	content = extractJson.Find(content)

	// Fix newlines inside strings
	// replace 2nd and 5th capturing group's char(10) newlines inside strings to `\n`
	content = stringExtraction.ReplaceAllFunc(content, func(s []byte) []byte {
		matches := stringExtraction.FindSubmatch(s)
		// replace newlines inside strings with \n
		matches[2] = bytes.ReplaceAll(matches[2], newLineFix, newLineReplacement)
		matches[5] = bytes.ReplaceAll(matches[5], newLineFix, newLineReplacement)
		return bytes.Join(matches[1:], nil)
	})

	// Fix escaped parentheses e.g. `\(text\)` to `\\(text\\)`
	content = fixParentheses.ReplaceAll(content, parenthesesReplacement)
	content = removeComments.ReplaceAll(content, nil)
	return content
}
