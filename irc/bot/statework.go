package bot

import (
	"log"
	"github.com/benishor/go-irc-bot/irc"
	"github.com/benishor/go-irc-bot/irc/replies"
)

type OnChannelStateHandler struct {
}

func (state *OnChannelStateHandler) HandleCommand(command *irc.IrcCommand, bot *Bot) {
	switch command.Command {
	case replies.RplJoin:
		user := command.Source;
		channel := command.Target;

		stats := bot.GetChannelStats(channel)
		stats.AddUser(irc.ParseIrcUser(user).Nickname, "")
		log.Println("Current users on channel: ", stats.Users)

		for _, v := range Plugins {
			v.OnUserJoin(bot, channel, user);
		}
	case replies.RplPart:
		user := command.Source;
		channel := command.Target;
		partMessage := command.Params;

		stats := bot.GetChannelStats(channel)
		stats.RemoveUser(irc.ParseIrcUser(user).Nickname)
		log.Println("Current users on channel: ", stats.Users)

		for _, v := range Plugins {
			v.OnUserPart(bot, channel, user, partMessage);
		}
	case replies.RplQuit:
		user := command.Source;
		quitMessage := command.Params;

		nickname := irc.ParseIrcUser(user).Nickname
		for _, stats := range bot.Channels {
			stats.RemoveUser(nickname)
			log.Println("Current users on channel: ", stats.Users)
		}

		for _, v := range Plugins {
			v.OnUserQuit(bot, user, quitMessage);
		}
	case replies.CmdPrivMsg:
		userMessage := command.Params
		issuer := command.Source;
		channel := command.Target;

		for _, v := range Plugins {
			v.OnChannelMessage(bot, channel, issuer, userMessage);
		}
	default:
		log.Printf("Unhandled command [%s]", command)
	}
}
