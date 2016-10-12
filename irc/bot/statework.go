package bot

import (
	"strings"
	"fmt"
	"log"
	"github.com/benishor/go-irc-bot/irc"
	"github.com/benishor/go-irc-bot/irc/commands"
	"github.com/benishor/go-irc-bot/irc/replies"
)

type OnChannelStateHandler struct {
}

func (state *OnChannelStateHandler) HandleCommand(command *irc.IrcCommand, bot *Bot) {
	switch command.Command {
	case replies.RplJoin:
		nicknameJoined := strings.Split(command.Source, "!")[0]
		bot.Settings.Output <- commands.SendMessage(command.Target, fmt.Sprintf("Hello there, dear %s!", nicknameJoined))
	default:
		log.Printf("Unhandled command [%s]", command)
	}
}

