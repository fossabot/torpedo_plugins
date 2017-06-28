package torpedo_file_plugin

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/tb0hdan/torpedo_plugins/torpedo_file_plugin/file"

	common "github.com/tb0hdan/torpedo_common"
	"github.com/tb0hdan/torpedo_registry"
)

func GetSetImageProcessMessage(api *torpedo_registry.BotAPI, channel_i interface{}, incoming_message string) {

	requestedFeature, command, message := common.GetRequestedFeature(incoming_message)
	channel := channel_i.(string)
	if command != "" {
		switch requestedFeature {
		case fmt.Sprintf("%sgetimg", api.CommandPrefix):
			fpath, mimetype, err := file.GetChannelFile(channel, command)
			if fpath != "" {
				common.ChannelsUploadImage([]string{channel}, command, fpath, mimetype, api)
				return
			} else {
				message = fmt.Sprintf("%+v", err)
			}
		case fmt.Sprintf("%ssetimg", api.CommandPrefix):
			msg, err := file.SetChannelFile(channel, command, api.CommandPrefix)
			if err != nil {
				message = fmt.Sprintf("%+v", err)
			} else {
				message = msg
			}
		case fmt.Sprintf("%slistimg", api.CommandPrefix), fmt.Sprintf("%slsimg", api.CommandPrefix):
			files, err := file.ListChannelFiles(channel)
			if err != nil {
				message = "An error occured while retrieving image file list"
			} else {
				message = ""
				for _, fname := range files {
					msg, err := base64.URLEncoding.DecodeString(fname)
					if err != nil {
						continue
					}
					message += fmt.Sprintf("`%s`\n", msg)
				}
				if message == "" {
					message = fmt.Sprintf("No files found, upload using %ssetimg first", api.CommandPrefix)
				} else {
					message = fmt.Sprintf("Available image files:\n%s", message)
				}
			}
		case fmt.Sprintf("%srmimg", api.CommandPrefix):
			fpath, _, _ := file.GetChannelFile(channel, command)
			if fpath != "" {
				err := os.Remove(fpath)
				if err != nil {
					message = fmt.Sprintf("An error occured while trying to remove `%s`", command)
				} else {
					message = fmt.Sprintf("Requested file `%s` was removed", command)
				}
			} else {
				message = fmt.Sprintf("Requested file `%s` was not found", command)
			}
		default:
			// should never get here
			message = "Unknown feature requested"
		}
	}
	api.Bot.PostMessage(channel, message, api)
}

func init() {
	torpedo_registry.Config.RegisterHandler("getimg", GetSetImageProcessMessage)
	torpedo_registry.Config.RegisterHelp("getimg", "Get stored image for this channel")
	torpedo_registry.Config.RegisterHandler("lsimg", GetSetImageProcessMessage)
	torpedo_registry.Config.RegisterHelp("lsimg", "List stored images for this channel")
	torpedo_registry.Config.RegisterHandler("listimg", GetSetImageProcessMessage)
	torpedo_registry.Config.RegisterHelp("listimg", "List stored images for this channel")
	torpedo_registry.Config.RegisterHandler("setimg", GetSetImageProcessMessage)
	torpedo_registry.Config.RegisterHelp("setimg", "Set stored image for this channel")
	torpedo_registry.Config.RegisterHandler("rmimg", GetSetImageProcessMessage)
	torpedo_registry.Config.RegisterHelp("rmimg", "Remove stored image for this channel")
}
