package commands

import "fmt"

func RegisterUser(username string, fullName string) (string) {
	return fmt.Sprintf("USER %s * * :%s\n", username, fullName)
}

func SetNickname(nickname string) (string) {
	return fmt.Sprintf("NICK %s\n", nickname)
}

func JoinChannel(channel string) (string) {
	return fmt.Sprintf("JOIN %s\n", channel)
}

func SendMessage(destination string, message string) (string) {
	return fmt.Sprintf("PRIVMSG %s :%s\n", destination, message)
}

func Pong(timestamp string) (string) {
	return fmt.Sprintf("PONG %s\n", timestamp)
}