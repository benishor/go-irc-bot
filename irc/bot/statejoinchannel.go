package bot

import (
	"log"
	"github.com/benishor/go-irc-bot/irc"
	"github.com/benishor/go-irc-bot/irc/commands"
	"github.com/benishor/go-irc-bot/irc/replies"
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
	case replies.RplJoin:
		bot.Write(commands.SendMessage(command.Target, "Hello all. Your obedient greeting slave is online"))
		bot.State = &OnChannelStateHandler{}
	default:
		log.Printf("Unhandled command [%s]", command)
	}
}


