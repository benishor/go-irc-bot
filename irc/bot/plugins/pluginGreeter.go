package plugins

import (
	"github.com/benishor/go-irc-bot/irc/bot"
	"github.com/benishor/go-irc-bot/irc/commands"
	"fmt"
)

func init() {
	bot.RegisterPlugin(&GreeterPlugin{})
}

type GreeterPlugin struct {
	bot.DefaultBotPlugin
}

func (p*GreeterPlugin) OnJoinChannel(bot *bot.Bot, channel string) {
	bot.Write(commands.SendMessage(channel, "Your greeting slave is online!"))
}

func (p*GreeterPlugin) OnUserJoin(bot *bot.Bot, channel string, user string) {
	bot.Write(commands.SendMessage(channel, fmt.Sprintf("Hi %s, may you have an enjoyable stay in %s!", user, channel)))
}
