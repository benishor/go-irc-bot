package main

import (
	"github.com/benishor/go-irc-bot/irc/bot"
	"github.com/benishor/go-irc-bot/irc/communication"
	"os"
	"os/signal"
	"syscall"
	_ "github.com/benishor/go-irc-bot/plugins" // self registering plugins
)

func main() {
	tcpChannel := communication.NewTcpChannel("irc.freenode.net", 6667)
	//tcpChannel := communication.NewTcpChannel("bucharest.ro.eu.undernet.org", 6667)

	defer tcpChannel.Close()

	config := &bot.Config{
		Nickname: "benishor",
		Channel: "#go-test-bot",
		FullName: "Evil g0 b0t"}

	robot := bot.NewBot(config, tcpChannel)

	registerShutdownHook(robot.Close)

	robot.Run()
}

func registerShutdownHook(closeHandler func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChannel {
			closeHandler()
		}
	}()
}