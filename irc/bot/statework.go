package bot

import (
	"strings"
	"fmt"
	"log"
	"github.com/benishor/go-irc-bot/irc"
	"github.com/benishor/go-irc-bot/irc/replies"
	"regexp"
	"github.com/benishor/go-irc-bot/irc/commands"
	"time"
)

type LastSeenInfo struct {
	User        string
	Channel     string
	LeavingTime time.Time
	Message     string
}

type LastSeenInfoMap map[string]LastSeenInfo

type LastSeenTool struct {
	info LastSeenInfoMap
}

func NewLastSeenTool() (*LastSeenTool) {
	return &LastSeenTool{
		info: make(map[string]LastSeenInfo)}
}

func (t*LastSeenTool) set(user, channel, message string) {
	nickname := strings.ToLower(strings.Split(user, "!")[0])
	t.info[nickname] = LastSeenInfo{
		User: user,
		Channel: channel,
		LeavingTime: time.Now(),
		Message: message}
}

func (t*LastSeenTool) get(nickname string) (*LastSeenInfo) {
	if val, ok := t.info[strings.ToLower(nickname)]; ok {
		return &val
	}
	return nil
}

// -----------------------------------------------------------------

type ChannelStats struct {
	Topic     string
	Nicknames [] string
}

type OnChannelStateHandler struct {
	nicknames []string
}

var botCommandPattern = regexp.MustCompile(`!([^\W]+)(.*)`)
var lastSeenTool = NewLastSeenTool()

func (state *OnChannelStateHandler) addUserToChannel(user, channel string) {
	nickname := strings.Split(user, "!")[0]
	state.nicknames = append(state.nicknames, nickname)
	log.Printf("User %s (%s) joined %s\n", nickname, user, channel)
}

func (state *OnChannelStateHandler) removeUserFromChannel(user, channel, reason string) {
	nickname := strings.Split(user, "!")[0]

	var i int
	for k, v := range state.nicknames {
		if v == nickname {
			i = k
			break;
		}
	}

	state.nicknames = append(state.nicknames[:i], state.nicknames[i + 1:]...)
	log.Printf("User %s (%s) left %s saying %s\n", nickname, user, channel, reason)

	lastSeenTool.set(user, channel, reason)
}

func (state *OnChannelStateHandler) HandleCommand(command *irc.IrcCommand, bot *Bot) {
	switch command.Command {
	case replies.RplJoin:
		state.addUserToChannel(command.Source, command.Target)
	//bot.Write(commands.SendMessage(command.Target, fmt.Sprintf("Hello there, dear %s!", nicknameJoined)))
	case replies.RplPart:
		state.removeUserFromChannel(command.Source, command.Target, command.Params)
	case replies.RplQuit:
		state.removeUserFromChannel(command.Source, command.Target, command.Params)
	case replies.RplNameReply:
		nicksInThisLine := strings.Split(command.Params, " ")
		for _, nick := range nicksInThisLine {
			state.nicknames = append(state.nicknames, nick)
		}
	case replies.RplEndOfNames:
		log.Println("Current users on channel: ", state.nicknames)
	case replies.CmdPrivMsg:
		userMessage := command.Params

		result := botCommandPattern.FindStringSubmatch(userMessage)
		if result != nil {
			botCommand := result[1]
			arguments := strings.TrimSpace(result[2])

			reply, err := state.dispatchBotCommand(command.Source, command.Target, botCommand, arguments)
			if err == nil {
				bot.Write(commands.SendMessage(command.Target, reply))
			}
		}
	default:
		log.Printf("Unhandled command [%s]", command)
	}
}

func (state *OnChannelStateHandler) dispatchBotCommand(issuer string, channel string, botCommand string, arguments string) (string, error) {
	log.Printf("Dispatching bot command [%s] with arguments [%s] issued by [%s] on [%s]",
		botCommand, arguments, issuer, channel)

	issuerNickname := strings.Split(issuer, "!")[0]
	if botCommand == "seen" {
		if arguments == "" {
			return "Correct command usage is \"!seen [nickname]\"", nil
		}
		userCurrentlyOnChannel := false
		for _, nick := range state.nicknames {
			if nick == arguments {
				userCurrentlyOnChannel = true
				break
			}
		}
		if userCurrentlyOnChannel {
			return fmt.Sprintf("Look closer, %s. %s is in the channel!", issuerNickname, arguments), nil
		}
		lastSeenInfo := lastSeenTool.get(arguments)
		if lastSeenInfo == nil {
			return fmt.Sprintf("%s: sorry, I didn't see %s", issuerNickname, arguments), nil
		} else {
			return fmt.Sprintf("%s (%s) was last seen on %s stating [%s]",
				arguments,
				lastSeenInfo.User,
				lastSeenInfo.LeavingTime.Format(time.UnixDate),
				lastSeenInfo.Message), nil
		}

	} else {
		return "", nil
	}
}

