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
	Model    string
	Models   []ModelData
	Device   string
}

// ModelData is part of DeviceData
type ModelData struct {
	Regex string
	Model string
}

// ParseClients parses a given YML file to an array of OssData
func ParseOss(file string) []OssData {
	content := getFileContent(file)
	ossData := make([]OssData, 0)
	err := yaml.Unmarshal(content, &ossData)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing yaml: %s", err.Error()))
	}
	for i, r := range ossData {
		reg := rewriteRegexp(r.Regex)
		ossData[i].Compiled = regexp.MustCompile(reg)
	}
	return ossData
}

// ParseClients parses a given YML file to an array of BotsData
func ParseBots(file string) []BotsData {
	content := getFileContent(file)
	botsData := make([]BotsData, 0)
	err := yaml.Unmarshal(content, &botsData)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing yaml: %s", err.Error()))
	}
	for i, r := range botsData {
		reg := rewriteRegexp(r.Regex)
		botsData[i].Compiled = regexp.MustCompile(reg)
	}
	return botsData
}

// ParseClients parses a given YML file to an array of ClientData and injects a type if it's not found in the yaml
func ParseClients(file, injectedType string) []ClientData {
	content := getFileContent(file)
	clientData := make([]ClientData, 0)
	err := yaml.Unmarshal(content, &clientData)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing yaml: %s", err.Error()))
	}
	for i, r := range clientData {
		reg := rewriteRegexp(r.Regex)
		clientData[i].Compiled = regexp.MustCompile(reg)
		if clientData[i].Type == "" {
			clientData[i].Type = injectedType
		}
	}
	return clientData
}

// ParseMultipleClient parses multiple client files assigning a type if not found in the yaml
func ParseMultipleClients(files map[string]string) []ClientData {
	clientData := make([]ClientData, 0)
	for file, injectedType := range files {
		clientData = append(clientData, ParseClients(file, injectedType)...)
	}
	return clientData
}

// ParseClients parses a given YML file to an array of DeviceData
func ParseDevice(file string) map[string]DeviceData {
	content := getFileContent(file)
	deviceData := make(map[string]DeviceData, 0)
	err := yaml.Unmarshal(content, &deviceData)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing yaml: %s", err.Error()))
	}
	for brand, r := range deviceData {
		reg := rewriteRegexp(r.Regex)
		tmp := deviceData[brand]
		tmp.Compiled = regexp.MustCompile(reg)
		deviceData[brand] = tmp
	}
	return deviceData
}

// ParseClients parses multiple YML files to an array of DeviceData
func ParseMultipleDevices(files []string) map[string]DeviceData {
	deviceData := make(map[string]DeviceData, 0)
	for _, file := range files {
		parsed := ParseDevice(file)
		for brand, data := range parsed {
			deviceData[brand] = data
		}
	}
	return deviceData
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
