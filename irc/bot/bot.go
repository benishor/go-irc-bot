package bot

import (
	"time"
	"github.com/benishor/go-irc-bot/irc"
	"github.com/benishor/go-irc-bot/irc/replies"
	"github.com/benishor/go-irc-bot/irc/commands"
	"log"
)

type Config struct {
	Nickname    string
	FullName    string
	Channel     string
	Input       chan string
	Output      chan string
	QuitChannel chan bool
}

type IrcBotStateHandler interface {
	HandleCommand(command *irc.IrcCommand, bot *Bot)
}

type Bot struct {
	State    IrcBotStateHandler
	Settings *Config
}

func NewBot(settings *Config) (*Bot) {
	return &Bot{
		Settings: settings,
		State : &RegistrationStateHandler{}}
}

func (bot *Bot)Run() {
	for true {
		select {
		case line := <-bot.Settings.Input:
			bot.handleLine(line)
		case <-bot.Settings.QuitChannel:
			log.Println("Received quit signal. Closing.")
			return;
		case <-time.After(time.Second * 1):
		//fmt.Print(".")
		}
	}

}

func (bot*Bot) handleLine(line string) {
	command, err := irc.ParseIrcCommand(line)
	if err == nil {
		switch command.Command {
		case replies.CmdPing:
			bot.Settings.Output <- commands.Pong(command.Params)
		default:
			bot.State.HandleCommand(command, bot)
		}
	}
}

