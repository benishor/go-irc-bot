package bot

import (
	"fmt"
	"strings"
	"log"
	"github.com/benishor/go-irc-bot/irc"
	"github.com/benishor/go-irc-bot/irc/commands"
	"github.com/benishor/go-irc-bot/irc/replies"
)

type RegistrationStateHandler struct {
	registrationSent bool
	nickTimesTaken   int
}

func (state *RegistrationStateHandler) HandleCommand(command *irc.IrcCommand, bot *Bot) {
	if !state.registrationSent {
		bot.Settings.Output <- commands.RegisterUser(bot.Settings.Nickname, bot.Settings.FullName)
		bot.Settings.Output <- commands.SetNickname(bot.Settings.Nickname)
		state.registrationSent = true
	}

	switch command.Command {
	case replies.ErrNicknameInUse:
		bot.Settings.Output <- commands.SetNickname(state.nextNickname(bot.Settings.Nickname))
	case replies.ErrErronoeusNickname:
		log.Fatal("ErronoeusNickname!")
	case replies.ErrNoMotd:
		fallthrough
	case replies.RplEndOfMotd:
		bot.State = &JoinChannelStateHandler{}
	default:
		log.Printf("Unhandled command [%s]", command)
	}
}

func (state *RegistrationStateHandler) nextNickname(nickname string) (string) {
	newNickname := fmt.Sprintf("%s%s", nickname, strings.Repeat("_", state.nickTimesTaken + 1))
	state.nickTimesTaken++
	return newNickname
}
