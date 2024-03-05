package utils

import (
	"regexp"
)

var (
	Patterns = map[string]*regexp.Regexp{
		"steps":      regexp.MustCompile(`(?i)steps:? (?P<steps>\d+)`),
		"sampler":    regexp.MustCompile(`(?i)sampler:? (?P<sampler>[\w+ ]+)`),
		"cfg":        regexp.MustCompile(`(?i)cfg scale:? (?P<cfg>[\d.]+)`),
		"seed":       regexp.MustCompile(`(?i)seed:? (?P<seed>\d+)`),
		"width":      regexp.MustCompile(`(?i)size:? (?P<width>\d+)x\d+`),
		"height":     regexp.MustCompile(`(?i)size:? \d+x(?P<height>\d+)`),
		"hash":       regexp.MustCompile(`(?i)model hash:? (?P<hash>\w+)`),
		"model":      regexp.MustCompile(`(?i)model ?[^h]:? (?P<model>[\w-]+)`),
		"denoising":  regexp.MustCompile(`(?i)denoising strength:? (?P<denoising>[\d.]+)`),
		"loraHashes": regexp.MustCompile(loraHashes),
		"tiHashes":   regexp.MustCompile(tiHashes),
		"version":    regexp.MustCompile(`(?i)version:? (?P<version>v[\w.-]+)`),
	}

	positivePattern = regexp.MustCompile(`(?is)(?:(?:primary |pos(?:itive)? )?prompts?:?)(.+)negative prompt:?`)
	negativePattern = regexp.MustCompile(`(?is)(?:(?:|neg(?:ative)? )?prompts?:?)(.*?)(?:steps|sampler|model)`)
	bbCode          = regexp.MustCompile(`\[\/?[\w=]+\]`)

	extractJson     = regexp.MustCompile(`(?ms){.*}`)
	removeComments  = regexp.MustCompile(`(?m)//.*$`)
	escapeBackslash = regexp.MustCompile(`\\+([()])`)
)

const (
	loraHashes = `(?i)lora hashes:? "(?P<lora>[^"]+)"`
	loraEntry  = `(?i)(?P<key>\w+): (?P<value>\w+)`
	tiHashes   = `(?i)ti hashes:? "(?P<ti>[^"]+)"`
	tiEntry    = `(?i)(?P<key>\w+): (?P<value>\w+)`

	escapeBackslashReplacement = `\\$1`
)

func RemoveBBCode(s string) string {
	return bbCode.ReplaceAllString(s, "")
}

func ExtractAll(s string, reg map[string]*regexp.Regexp) map[string]string {
	result := make(map[string]string)

	for key, r := range reg {
		result[key] = Extract(s, r)
	}

	return result
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

func ExtractJson(content string) string {
	content = extractJson.FindString(content)
	content = removeComments.ReplaceAllString(content, "")
	content = escapeBackslash.ReplaceAllString(content, escapeBackslashReplacement)
	return content
}
