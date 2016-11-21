package plugins

import (
	"github.com/benishor/go-irc-bot/irc/bot"
	"strings"
	"github.com/benishor/go-irc-bot/irc/commands"
	"github.com/mvdan/xurls"
	"net/http"
	"log"
	"golang.org/x/net/html"
)

func init() {
	bot.RegisterPlugin(&UrlTitlePlugin{})
}

type UrlTitlePlugin struct {
	bot.DefaultBotPlugin
}

func (p*UrlTitlePlugin) OnChannelMessage(bot *bot.Bot, channel string, user string, message string) {
	titles := []string{}
	for _, url := range extractUrlsFromMessage(message) {
		title := getTitleForUrl(url)
		if title != "" {
			titles = append(titles, title)
		}
	}

	if len(titles) > 0 {
		bot.Write(commands.SendMessage(channel, strings.Join(titles, ", ")))
	}
}

func extractUrlsFromMessage(message string) []string {
	return xurls.Relaxed.FindAllString(message, 10)
}

func getTitleForUrl(url string) string {
	response, err := http.Head(url)
	if err != nil {
		return ""
	}

	if response.ContentLength > 1 * 1024 * 1024 {
		log.Printf("Dropping request for url %s due to large size: %l", url, response.ContentLength)
		return ""
	}

	response, err = http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch url %s. Cause: %s", url, err)
		return ""
	}

	// parse content
	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Printf("Failed to parse content from %s. Cause: %s", url, err)
		return ""
	}

	var title string = ""

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				title = c.Data
				break;
			}
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(doc)
	return strings.TrimSpace(title);
}
