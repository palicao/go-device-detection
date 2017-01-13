package lib

import (
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
	Engine   DetectionRegexEngine
}

// DetectionRegexEngine is a part of the regex representing the engine information
type DetectionRegexEngine struct {
	Default  string
	Versions map[string]string
}

// ParseYml parses a given YML file to an array of BrowserRegex
func ParseYml(file string) []DetectionRegex {
	content := getFileContent(file)
	regexps := make([]DetectionRegex, 0)
	yaml.Unmarshal(content, &regexps)
	for i, r := range regexps {
		if !strings.Contains(r.Regex, "?!") {
			regexps[i].Compiled = regexp.MustCompile(strings.Replace(r.Regex, "++", "{1,}", -1))
		}
	}
	return regexps
}

func getFileContent(file string) []byte {
	f, err := ioutil.ReadFile(file)
	if err == nil {
		return f
	}
	return make([]byte, 0)
}
