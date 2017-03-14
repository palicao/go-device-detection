package lib

import (
	"errors"
	"strconv"
	"strings"
)

var ossRegexps = ParseYml("./piwik/regexes/oss.yml")
var botsRegexps = ParseYml("./piwik/regexes/bots.yml")
var clientRegexps = buildClientRegexps()
//var deviceRegexps = buildDeviceRegexps()

func buildClientRegexps() []DetectionRegex {
	var feedReader = ParseYml("./piwik/regexes/client/feed_readers.yml")
	var mobileApps = InjectType(ParseYml("./piwik/regexes/client/mobile_apps.yml"), "Mobile App")
	var mediaPlayers = InjectType(ParseYml("./piwik/regexes/client/mediaplayers.yml"), "Mediaplayer")
	var pim = InjectType(ParseYml("./piwik/regexes/client/pim.yml"), "PIM")
	var browsers = InjectType(ParseYml("./piwik/regexes/client/browsers.yml"), "Browser")
	var libraries = InjectType(ParseYml("./piwik/regexes/client/libraries.yml"), "Library")

	total := make([]DetectionRegex, 0)

	total = append(total, feedReader...)
	total = append(total, mobileApps...)
	total = append(total, mediaPlayers...)
	total = append(total, pim...)
	total = append(total, browsers...)
	total = append(total, libraries...)

	return total
}

type DetectedDevice interface {
	IsBot() bool
	IsBrowser() bool
	IsFeedReader() bool
	IsMobileApp() bool
	IsPIM() bool
	IsLibrary() bool
	IsMediaPlayer() bool
}

type DetectedDeviceInfo struct {
	UserAgent  string
	BotInfo    BotInfo
	OSInfo     OSInfo
	DeviceInfo DeviceInfo
	ClientInfo ClientInfo
	isBot      bool
}

type OSInfo struct {
	Name    string
	Version string
}

type BotInfo struct {
	Name         string
	Category     string
	Url          string
	ProducerName string
	ProducerUrl  string
}

type DeviceInfo struct {
	Name  string
	Type  string
	Model string
}

type ClientInfo struct {
	Name    string
	Type    string
	Version string
	Engine  string
}

func (d *DetectedDeviceInfo) IsBot() bool {
	return d.isBot == true
}

func (d *DetectedDeviceInfo) IsBrowser() bool {
	return d.ClientInfo.Type == "Browser"
}

func (d *DetectedDeviceInfo) IsFeedReader() bool {
	return strings.Contains(d.ClientInfo.Type, "Feed Reader")
}

func (d *DetectedDeviceInfo) IsMobileApp() bool {
	return d.ClientInfo.Type == "Mobile App"
}

func (d *DetectedDeviceInfo) IsPIM() bool {
	return d.ClientInfo.Type == "PIM"
}

func (d *DetectedDeviceInfo) IsLibrary() bool {
	return d.ClientInfo.Type == "Library"
}

func (d *DetectedDeviceInfo) IsMediaPlayer() bool {
	return d.ClientInfo.Type == "Mediaplayer"
}

func Detect(ua string) (DetectedDeviceInfo, error) {
	device := DetectedDeviceInfo{}
	device.UserAgent = ua

	botInfo, err := detectBot(ua)
	if err == nil {
		device.isBot = true
		device.BotInfo = botInfo
		return device, nil
	}

	osInfo, err := detectOS(ua)
	if err == nil {
		device.OSInfo = osInfo
	}

	clientInfo, err := detectClient(ua)
	if err == nil {
		device.ClientInfo = clientInfo
	}

	return device, errors.New("bah")
}

func detectBot(ua string) (BotInfo, error) {
	botInfo := BotInfo{}
	for _, bot := range botsRegexps {
		found := bot.Compiled.FindStringSubmatch(ua)
		if len(found) > 0 {
			botInfo.Name = bot.Name
			botInfo.Category = bot.Category
			botInfo.Url = bot.Url
			botInfo.ProducerName = bot.Producer.Name
			botInfo.ProducerUrl = bot.Producer.Url

			return botInfo, nil
		}
	}
	return botInfo, errors.New("Not a bot")
}

func detectOS(ua string) (OSInfo, error) {
	osInfo := OSInfo{}
	for _, oss := range ossRegexps {
		found := oss.Compiled.FindStringSubmatch(ua)
		if len(found) > 0 {
			osInfo.Name = oss.Name
			osInfo.Version = parseVersion(oss.Version, found)

			return osInfo, nil
		}
	}
	return osInfo, errors.New("Unknown OS")
}

func detectClient(ua string) (ClientInfo, error) {
	clientInfo := ClientInfo{}
	for _, client := range clientRegexps {
		found := client.Compiled.FindStringSubmatch(ua)
		if len(found) > 0 {
			clientInfo.Name = client.Name
			clientInfo.Version = parseVersion(client.Version, found)
			clientInfo.Type = client.Type

			return clientInfo, nil
		}
	}
	return clientInfo, errors.New("Unknown client")
}

func parseVersion(matcher string, version []string) string {
	if matcher[0:1] == "$" {
		part, _ := strconv.Atoi(matcher[1:])
		return version[part]
	}
	return matcher
}
