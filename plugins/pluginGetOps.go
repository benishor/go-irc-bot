package plugins

import (
	"github.com/benishor/go-irc-bot/irc/bot"
	"github.com/benishor/go-irc-bot/irc/commands"
	"time"
	"log"
	"strings"
)

func init() {
	bot.RegisterPlugin(&GetOpsPlugin{})
}

type GetOpsPlugin struct {
	bot.DefaultBotPlugin
}

func (p*GetOpsPlugin) OnUserPart(bot *bot.Bot, channel string, user string, partMessage string) {
	log.Printf("OnUserPart from %s", channel)

	justMeOnTheChannel := len(bot.GetChannelStats(channel).Users) == 1

	if justMeOnTheChannel && !strings.Contains(bot.GetChannelStats(channel).Users[0].Modes, "@") {
		log.Printf("Hopping in order to get ops on %s", channel)
		bot.Write(commands.PartChannel(channel))
		go func() {
			time.Sleep(time.Duration(1) * time.Second)
			bot.Write(commands.JoinChannel(channel))
		}()
	}
}
