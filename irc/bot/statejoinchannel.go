package bot

import (
	"log"
	"github.com/benishor/go-irc-bot/irc"
	"github.com/benishor/go-irc-bot/irc/commands"
	"github.com/benishor/go-irc-bot/irc/replies"
	"strings"
)

type JoinChannelStateHandler struct {
	joinSent bool
}

func (state *JoinChannelStateHandler) HandleCommand(command *irc.IrcCommand, bot *Bot) {
	if !state.joinSent {
		bot.Write(commands.JoinChannel(bot.Config.Channel))
		state.joinSent = true
	}

	switch command.Command {
	case replies.RplTopic:
		targetPieces := strings.Split(command.Target, " ")
		channel := targetPieces[len(targetPieces) - 1]
		bot.GetChannelStats(channel).Topic = command.Params
		log.Println("Got channel topic");
	case replies.RplNoTopic:
		log.Println("No topic");
	case replies.RplJoin:
		log.Println("Got join message ");
		//bot.Write(commands.SendMessage(command.Target, "Hello all. Your obedient greeting slave is online"))
	case replies.RplNameReply:
		targetChannel := irc.ExtractChannelFromNamesReply(command.Target)
		stats := bot.GetChannelStats(targetChannel)

		namesInChannel := strings.Split(command.Params, " ")
		for _, nickWithModes := range namesInChannel {
			nickname, modes := irc.ParseNickWithModes(nickWithModes, bot.ServerSettings.Prefixes)
			stats.AddUser(nickname, modes);
		}

	case replies.RplEndOfNames:
		targetChannel := irc.ExtractChannelFromNamesReply(command.Target)
		log.Println("Current users on channel: ", bot.Channels[targetChannel].Users)

		for _, v := range Plugins {
			v.OnJoinChannel(bot, targetChannel)
		}
		// we really joined the channel. change state
		bot.State = &OnChannelStateHandler{}
	default:
		log.Printf("Unhandled command [%s]", command)
	}
}


