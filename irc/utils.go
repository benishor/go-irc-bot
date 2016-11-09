package irc

import (
	"regexp"
	"fmt"
	"strings"
)

type IrcCommand struct {
	Source  string
	Target  string
	Command string
	Params  string
}

type IrcUser struct {
	Nickname string
	User     string
	Host     string
}

func (cmd*IrcCommand)String() string {
	return fmt.Sprintf("IrcCommand{Command:\"%s\", Source:\"%s\", Target:\"%s\", Params:\"%s\"}",
		cmd.Command,
		cmd.Source,
		cmd.Target,
		cmd.Params)
}

var commandPattern = regexp.MustCompile(`^(?:[:](\S+) )?(\S+)(?: ([^:].+?))?(?: [:](.+))?$`)

func ParseIrcCommand(line string) (*IrcCommand, error) {
	if commandPattern.MatchString(line) {
		matches := commandPattern.FindStringSubmatch(line)

		return &IrcCommand{
			Source: matches[1],
			Target: matches[3],
			Command: matches[2],
			Params: matches[4]}, nil
	}
	return nil, fmt.Errorf("Failed to parse command from [%s]", line)
}

func ExtractChannelFromNamesReply(target string) string {
	channelPieces := strings.Split(target, " ");
	return strings.TrimSpace(channelPieces[len(channelPieces) - 1])
}

func ParseIrcUser(user string) (IrcUser) {
	nickAndMask := strings.Split(user, "!")
	userAndHost := strings.Split(nickAndMask[1], "@")
	return IrcUser{
		Nickname: nickAndMask[0],
		User: userAndHost[0],
		Host: userAndHost[1]}
}

func ParseNickWithModes(nickWithModes, serverPrefixes string) (nickname, modes string) {
	modeCharacters := 0
	for _, char := range nickWithModes {
		if strings.ContainsRune(serverPrefixes, char) {
			modeCharacters++;
		} else {
			break;
		}
	}
	userChannelModes := nickWithModes[0:modeCharacters]
	nickWithoutModes := nickWithModes[modeCharacters:]
	return nickWithoutModes, userChannelModes
}