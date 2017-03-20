package lib

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// OssData is needed for operating systems detection
type OssData struct {
	Regex    string
	Compiled *regexp.Regexp
	Name     string
	Version  string
}

// BotsData is needed for bot detection
type BotsData struct {
	Regex    string
	Compiled *regexp.Regexp
	Name     string
	Category string
	Url      string
	Producer ProducerData
}

// ProducerData is part of the regex representing bot producer data
type ProducerData struct {
	Name string
	Url  string
}

// ClientData is for browsers and other clients detection
type ClientData struct {
	Regex    string
	Compiled *regexp.Regexp
	Name     string
	Version  string
	Type     string
	Engine   EngineData
}

// EngineData is a part of the regex representing the engine information (for browsers)
type EngineData struct {
	Default  string
	Versions map[string]string
}

// DeviceData is needed for device detection
type DeviceData struct {
	Regex    string
	Compiled *regexp.Regexp
	Model  string
	Models []ModelData
	Device string
}

// ModelData is part of DeviceData
type ModelData struct {
	Regex string
	Model string
}

// ParseYml parses a given YML file to an array of BrowserRegex
func ParseYml(file string) []ClientData {
	content := getFileContent(file)
	regexps := make([]ClientData, 0)
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

func InjectType(clientRegexes []ClientData, t string) []ClientData {
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
