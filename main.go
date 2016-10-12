package main

import (
	"github.com/benishor/go-irc-bot/irc/bot"
	"github.com/benishor/go-irc-bot/irc/communication"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	readChannel := make(chan string, 10)
	writeChannel := make(chan string, 10)
	quitChannel := make(chan bool, 1)

	tcpChannel := communication.NewChannel("irc.freenode.net", 6667)
	defer tcpChannel.Close()
	tcpChannel.BindIOChannels(readChannel, writeChannel)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for _ = range signalChannel {
			quitChannel <- true
		}
	}()

	settings := &bot.Config{
		Nickname: "benishor",
		Input: readChannel,
		Output: writeChannel,
		QuitChannel: quitChannel,
		Channel: "#go-test-bot",
		FullName: "Evil g0 b0t"}

	robot := bot.NewBot(settings)
	robot.Run()
}