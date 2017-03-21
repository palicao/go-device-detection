package lib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var ossRegexps = ParseOss("./piwik/regexes/oss.yml")

var botsRegexps = ParseBots("./piwik/regexes/bots.yml")

var clientRegexps = ParseMultipleClients(map[string]string{
	"./piwik/regexes/client/feed_readers.yml": "",
	"./piwik/regexes/client/mobile_apps.yml":  "Mobile App",
	"./piwik/regexes/client/mediaplayers.yml": "Mediaplayer",
	"./piwik/regexes/client/pim.yml":          "PIM",
	"./piwik/regexes/client/browsers.yml":     "Browser",
	"./piwik/regexes/client/libraries.yml":    "Library",
})

var tvRegexps = ParseDevice("./piwik/regexes/device/televisions.yml")

var deviceRegexps = ParseMultipleDevices([]string{
	"./piwik/regexes/device/consoles.yml",
	"./piwik/regexes/device/car_browsers.yml",
	"./piwik/regexes/device/cameras.yml",
	"./piwik/regexes/device/portable_media_player.yml",
	"./piwik/regexes/device/mobiles.yml",
})

type Detected interface {
	IsBot() bool
	IsBrowser() bool
	IsFeedReader() bool
	IsMobileApp() bool
	IsPIM() bool
	IsLibrary() bool
	IsMediaPlayer() bool
}

type DetectionInfo struct {
	UserAgent  string
	BotInfo    BotInfo
	OSInfo     OSInfo
	DeviceInfo DeviceInfo
	ClientInfo ClientInfo
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
	Device string
	Model  string
}

type ClientInfo struct {
	Name    string
	Type    string
	Version string
	Engine  string
}

func (d *DetectionInfo) IsBot() bool {
	return d.BotInfo.Name != ""
}

func (d *DetectionInfo) IsBrowser() bool {
	return d.ClientInfo.Type == "Browser"
}

func (d *DetectionInfo) IsFeedReader() bool {
	return strings.Contains(d.ClientInfo.Type, "Feed Reader")
}

func (d *DetectionInfo) IsMobileApp() bool {
	return d.ClientInfo.Type == "Mobile App"
}

func (d *DetectionInfo) IsPIM() bool {
	return d.ClientInfo.Type == "PIM"
}

func (d *DetectionInfo) IsLibrary() bool {
	return d.ClientInfo.Type == "Library"
}

func (d *DetectionInfo) IsMediaPlayer() bool {
	return d.ClientInfo.Type == "Mediaplayer"
}

func Detect(ua string) (DetectionInfo, error) {
	info := DetectionInfo{}
	info.UserAgent = ua

	botInfo, err := detectBot(ua)
	if err == nil {
		info.BotInfo = botInfo
		return info, nil
	}

	osInfo, err := detectOS(ua)
	if err == nil {
		info.OSInfo = osInfo
	}

	clientInfo, err := detectClient(ua)
	if err == nil {
		info.ClientInfo = clientInfo
	}

	deviceInfo, err := detectDevice(ua)
	if err == nil {
		info.DeviceInfo = deviceInfo
	}

	return info, errors.New(fmt.Sprintf("Couldn't detect User Agent %s", ua))
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

func detectDevice(ua string) (DeviceInfo, error) {
	hbbRegex := regexp.MustCompile("HbbTV/([1-9]{1}(?:\\.[0-9]{1}){1,2})")
	isHbb := hbbRegex.MatchString(ua)
	if isHbb {
		return detectDeviceBetween(ua, tvRegexps)
	}
	return detectDeviceBetween(ua, deviceRegexps)
}

func detectDeviceBetween(ua string, r map[string]DeviceData) (DeviceInfo, error) {
	deviceInfo := DeviceInfo{}
	for _, device := range r {
		found := device.Compiled.FindStringSubmatch(ua)
		if len(found) > 0 {
			deviceInfo.Device = device.Device
			deviceInfo.Model = device.Model
			return deviceInfo, nil
		}
	}
	return deviceInfo, errors.New("Unable to detect device")
}

func parseVersion(matcher string, version []string) string {
	if matcher[0:1] == "$" {
		part, _ := strconv.Atoi(matcher[1:])
		return strings.Replace(version[part], "_", ".", -1)
	}
	return matcher
}
