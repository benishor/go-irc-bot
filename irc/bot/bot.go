package bot

import (
	"time"
	"github.com/benishor/go-irc-bot/irc"
	"github.com/benishor/go-irc-bot/irc/replies"
	"github.com/benishor/go-irc-bot/irc/commands"
	"log"
	"github.com/benishor/go-irc-bot/irc/communication"
	"fmt"
)

type Config struct {
	Nickname string
	FullName string
	Channel  string
	Server   string
}

type ChannelUser struct {
	Nickname string
	Modes    string
}

type ChannelStats struct {
	Topic string
	Modes string
	Users []ChannelUser
}

func (user *ChannelUser) String() (string) {
	return fmt.Sprintf("%s%s", user.Modes, user.Nickname);
}

func (stats *ChannelStats) AddUser(username, modes string) {
	stats.Users = append(stats.Users, ChannelUser{
		Nickname: username,
		Modes: modes})
}

func (stats *ChannelStats) RemoveUser(nickname string) {
	for k, channelUser := range stats.Users {
		if channelUser.Nickname == nickname {
			stats.Users = append(stats.Users[:k], stats.Users[k + 1:]...)
			break;
		}
	}
}

type IrcBotStateHandler interface {
	HandleCommand(command *irc.IrcCommand, bot *Bot)
}

type ServerConfig struct {
	Prefixes string
}

type Bot struct {
	State                IrcBotStateHandler
	CurrentNickname      string
	Config               *Config
	Channels             map[string]*ChannelStats
	ServerSettings       ServerConfig
	quitChannel          chan bool
	readChannel          chan string
	writeChannel         chan string
	communicationChannel communication.Channel
}

func NewBot(config *Config, communicationChannel communication.Channel) (*Bot) {
	result := &Bot{
		Config: config,
		State: &RegistrationStateHandler{},
		Channels: make(map[string]*ChannelStats),
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
	go func() {
		time.Sleep(time.Duration(1) * time.Second)
		bot.quitChannel <- true
	}()
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

func (bot*Bot) GetChannelStats(channel string) (*ChannelStats) {
	_, ok := bot.Channels[channel];
	if !ok {
		bot.Channels[channel] = new(ChannelStats)
	}
	return bot.Channels[channel];
}
