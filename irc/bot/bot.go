package bot

import (
	"time"
	"github.com/benishor/go-irc-bot/irc"
	"github.com/benishor/go-irc-bot/irc/replies"
	"github.com/benishor/go-irc-bot/irc/commands"
	"log"
	"github.com/benishor/go-irc-bot/irc/communication"
)

type Config struct {
	Nickname string
	FullName string
	Channel  string
	Server   string
}

type IrcBotStateHandler interface {
	HandleCommand(command *irc.IrcCommand, bot *Bot)
}

type Bot struct {
	State                IrcBotStateHandler
	Config               *Config
	quitChannel          chan bool
	readChannel          chan string
	writeChannel         chan string
	communicationChannel communication.Channel
}

func NewBot(config *Config, communicationChannel communication.Channel) (*Bot) {
	result := &Bot{
		Config: config,
		State : &RegistrationStateHandler{},
		quitChannel : make(chan bool, 1),
		readChannel: make(chan string, 10),
		writeChannel: make(chan string, 10),
		communicationChannel: communicationChannel}

	result.bindIOChannels()
	return result
}

func (bot *Bot) bindIOChannels() {
	// constantly read from socket and put into input go channel
	go func(readChannel chan string, pipe communication.Channel) {
		for true {
			line, err := pipe.ReadLine()
			if err != nil {
				log.Fatalf("Error when reading. Reason: [%s]", err)
			} else {
				readChannel <- line
			}
		}
	}(bot.readChannel, bot.communicationChannel)

	// constantly read from output go channel and write on socket
	go func(channelToReadFrom chan string, pipe communication.Channel) {
		for true {
			pipe.WriteLine(<-channelToReadFrom)
		}
	}(bot.writeChannel, bot.communicationChannel)

}

func (bot *Bot) Run() {
	for true {
		select {
		case line := <-bot.readChannel:
			bot.handleLine(line)
		case <-bot.quitChannel:
			log.Println("Received quit signal. Closing.")
			return;
		case <-time.After(time.Second * 1):
		//fmt.Print(".")
		}
	}

}

func (bot *Bot) Close() {
	bot.Write("QUIT interrupted\n")
	bot.quitChannel <- true
}

func (bot *Bot) Write(content string) {
	bot.writeChannel <- content
}

func (bot*Bot) handleLine(line string) {
	command, err := irc.ParseIrcCommand(line)
	if err == nil {
		switch command.Command {
		case replies.CmdPing:
			bot.Write(commands.Pong(command.Params))
		default:
			bot.State.HandleCommand(command, bot)
		}
	}
}

