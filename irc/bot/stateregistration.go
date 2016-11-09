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
		bot.CurrentNickname = bot.Config.Nickname
		bot.Write(commands.RegisterUser(bot.Config.Nickname, bot.Config.FullName))
		bot.Write(commands.SetNickname(bot.CurrentNickname))
		state.registrationSent = true
	}

	switch command.Command {
	case replies.RplBounce:
		capabilities := strings.Split(command.Target, " ")
		for _, cap := range capabilities {
			//log.Println("Server capability: " + cap)
			if strings.HasPrefix(cap, "PREFIX=") {
				pieces := strings.Split(cap, "=")
				pieces2 := strings.Split(pieces[1], ")")
				bot.ServerSettings.Prefixes = pieces2[1]
			}
		}
	case replies.ErrNicknameInUse:
		bot.CurrentNickname = state.nextNickname(bot.Config.Nickname)
		bot.Write(commands.SetNickname(bot.CurrentNickname))
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
