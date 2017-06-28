package torpedo_shorteners_plugin

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/tb0hdan/torpedo_registry"
	common "github.com/tb0hdan/torpedo_common"
)

func QREncoderProcessMessage(api *torpedo_registry.BotAPI, channel interface{}, incoming_message string) {
	cu := &common.Utils{}
	cu.SetLoggerPrefix("shorteners-plugin")
	command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%sqr", api.CommandPrefix)))

	if command == "" {
		api.Bot.PostMessage(channel, fmt.Sprintf("Usage: %sqr query\n", api.CommandPrefix), api)
	} else {
		command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%sqr", api.CommandPrefix)))
		richmsg := torpedo_registry.RichMessage{ImageURL: fmt.Sprintf("http://chart.apis.google.com/chart?cht=qr&chs=350x350&chld=M|2&chl=%s", command), Text: command}
		api.Bot.PostMessage(channel, "", api, richmsg)
	}
}

func TinyURLProcessMessage(api *torpedo_registry.BotAPI, channel interface{}, incoming_message string) {
	cu := &common.Utils{}
	cu.SetLoggerPrefix("shorteners-plugin")
	command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%stinyurl", api.CommandPrefix)))

	if command == "" {
		api.Bot.PostMessage(channel, fmt.Sprintf("Usage: %stinyurl url\n", api.CommandPrefix), api)
	} else {
		command := strings.TrimSpace(strings.TrimLeft(incoming_message, fmt.Sprintf("%stinyurl", api.CommandPrefix)))
		query := url.QueryEscape(command)
		result, err := cu.GetURLBytes(fmt.Sprintf("https://tinyurl.com/api-create.php?url=%s", query))
		message := "An error occured during TinyURL encoding process"
		if err == nil {
			message = string(result)
		}
		api.Bot.PostMessage(channel, message, api)
	}
}

func CryptoProcessMessage(api *torpedo_registry.BotAPI, channel interface{}, incoming_message string) {
	requestedFeature, command, message := common.GetRequestedFeature(incoming_message)
	if command != "" {
		switch requestedFeature {
		case fmt.Sprintf("%sb64e", api.CommandPrefix):
			message = base64.StdEncoding.EncodeToString([]byte(command))
		case fmt.Sprintf("%sb64d", api.CommandPrefix):
			decoded, err := base64.StdEncoding.DecodeString(command)
			if err != nil {
				message = fmt.Sprintf("%v", err)
			} else {
				message = string(decoded)
			}
		case fmt.Sprintf("%smd5", api.CommandPrefix):
			message = common.MD5Hash(command)
		case fmt.Sprintf("%ssha1", api.CommandPrefix):
			message = common.SHA1Hash(command)
		case fmt.Sprintf("%ssha256", api.CommandPrefix):
			message = common.SHA256Hash(command)
		case fmt.Sprintf("%ssha512", api.CommandPrefix):
			message = common.SHA512Hash(command)
		default:
			// should never get here
			message = "Unknown feature requested"
		}
	}
	api.Bot.PostMessage(channel, message, api)
}

func init() {
	torpedo_registry.Config.RegisterHandler("qr", QREncoderProcessMessage)
	torpedo_registry.Config.RegisterHelp("qr", "Create QR Code from URL")
	torpedo_registry.Config.RegisterHandler("tinyurl", TinyURLProcessMessage)
	torpedo_registry.Config.RegisterHelp("tinyurl",  "Shorten URL using TinyURL.com")
	torpedo_registry.Config.RegisterHandler("b64e", CryptoProcessMessage)
	torpedo_registry.Config.RegisterHelp("b64e",  "Base64 encode")
	torpedo_registry.Config.RegisterHandler("b64d", CryptoProcessMessage)
	torpedo_registry.Config.RegisterHelp("b64d",  "Base64 decode")
	torpedo_registry.Config.RegisterHandler("md5", CryptoProcessMessage)
	torpedo_registry.Config.RegisterHelp("md5",  "Calculate message MD5 sum")
	torpedo_registry.Config.RegisterHandler("sha1", CryptoProcessMessage)
	torpedo_registry.Config.RegisterHelp("sha1", "Calculate message SHA1 sum")
	torpedo_registry.Config.RegisterHandler("sha256", CryptoProcessMessage)
	torpedo_registry.Config.RegisterHelp("sha256", "Calculate message SHA256 sum")
	torpedo_registry.Config.RegisterHandler("sha512", CryptoProcessMessage)
	torpedo_registry.Config.RegisterHelp("sha512", "Calculate message SHA512 sum")
}
