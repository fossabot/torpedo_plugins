package torpedo_bashorg_plugin

import (
	"golang.org/x/net/html"

	"fmt"
	"github.com/tb0hdan/torpedo_bashorg_plugin/bashorg"
	"github.com/tb0hdan/torpedo_registry"
)

func BashOrgProcessMessage(api *torpedo_registry.BotAPI, channel interface{}, incoming_message string) {
	item := api.Bot.GetCachedItem("bashorg")
	if item != "" {
		api.Bot.PostMessage(channel, item, api)
		return
	}

	bashorg := bashorg.NewClient()
	r := bashorg.Get_bashorg_html("http://bash.org/?random")
	quotes := make(map[int]string)

	i := 0

	var f func(*html.Node)
	var quote string

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" && len(n.Attr) > 0 {
			attr := n.Attr[0]
			if attr.Key == "class" && attr.Val == "qt" {
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

	quote = api.Bot.SetCachedItems("bashorg", quotes)
	api.Bot.PostMessage(channel, quote, api)
}

func init() {
	torpedo_registry.Config.RegisterHandler("bashorg", BashOrgProcessMessage)
	torpedo_registry.Config.RegisterHelp("bashorg", "Get http://bash.org random quote")
}
