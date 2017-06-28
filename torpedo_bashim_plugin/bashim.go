package torpedo_bashim_plugin

import (
	"golang.org/x/net/html"

	"fmt"
	"github.com/tb0hdan/torpedo_bashim_plugin/bashim"
	"github.com/tb0hdan/torpedo_registry"
)

func BashProcessMessage(api *torpedo_registry.BotAPI, channel interface{}, incoming_message string) {
	item := api.Bot.GetCachedItem("bashim")
	if item != "" {
		api.Bot.PostMessage(channel, item, api)
		return
	}

	bashim := bashim.NewClient()
	r := bashim.Get_html("http://bash.im/random")
	quotes := make(map[int]string)

	i := 0

	var f func(*html.Node)
	var quote string

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			attr := n.Attr[0]
			if attr.Key == "class" && attr.Val == "text" {
				quote = ""
				for mc := n.FirstChild; mc != nil; mc = mc.NextSibling {
					if n.Type == html.ElementNode && mc.Data == "br" {
						continue
					}
					quote += fmt.Sprintf("%s\n", mc.Data)
				}
				quotes[i] = quote
				i += 1
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(r)

	quote = api.Bot.SetCachedItems("bashim", quotes)
	api.Bot.PostMessage(channel, quote, api)
}

func init() {
	torpedo_registry.Config.RegisterHandler("bashim", BashProcessMessage)
	torpedo_registry.Config.RegisterHelp("bashim", "Get http://bash.im random quote")
}
