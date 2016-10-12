package irc

import (
	"regexp"
	"fmt"
)

type IrcCommand struct {
	Source  string
	Target  string
	Command string
	Params  string
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


