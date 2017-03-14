package lib

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// DetectionRegex contains all the parsed data from the yml file
type DetectionRegex struct {
	Regex    string
	Compiled *regexp.Regexp
	Name     string
	Version  string
	Category string
	Type     string
	Url      string
	Engine   EngineData
	Producer ProducerData
	Model    string
	Models   []ModelDetection
}

type ModelDetection struct {
	Regex string
	Model string
}

// EngineData is a part of the regex representing the engine information
type EngineData struct {
	Default  string
	Versions map[string]string
}

// ProducerData is part of the regex representing bot producer data
type ProducerData struct {
	Name string
	Url  string
}

// ParseYml parses a given YML file to an array of BrowserRegex
func ParseYml(file string) []DetectionRegex {
	content := getFileContent(file)
	regexps := make([]DetectionRegex, 0)
	err := yaml.Unmarshal(content, &regexps)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing yaml: %s", err.Error()))
	}
	for i, r := range regexps {
		reg := rewriteRegexp(r.Regex)
		regexps[i].Compiled = regexp.MustCompile(reg)
	}
	return regexps
}

func InjectType(clientRegexes []DetectionRegex, t string) []DetectionRegex {
	for i := 0; i < len(clientRegexes); i++ {
		clientRegexes[i].Type = t
	}
	return clientRegexes
}

// remove perl-specific regexp bits
func rewriteRegexp(reg string) string {
	r := reg

	// replace nested repetition operator ++ syntax with {1,}
	r = strings.Replace(r, "++", "{1,}", -1)

	// remove negative lookaheads and lookbehinds ("(?!" and "(?<")
	// @TODO find a working solution!
	negativeLookahead := regexp.MustCompile("\\(\\?[!|<][^)]*\\)")
	r = negativeLookahead.ReplaceAllString(r, "")

	return r
}

func getFileContent(file string) []byte {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintf("Error while reading file %s: %s", file, err))
	}
	return f
}
