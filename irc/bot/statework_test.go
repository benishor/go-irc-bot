package bot

import (
	"testing"
	"github.com/benishor/go-irc-bot/irc"
	"fmt"
)

func TestSeenCommand(t *testing.T) {
	bot := &Bot{
	}

	handler := OnChannelStateHandler{}
	command, _ := irc.ParseIrcCommand(":benishor_!~benny@86.122.19.112 PRIVMSG #undernet :!seen johnny")
	fmt.Println(command)
	handler.HandleCommand(command, bot);
}
