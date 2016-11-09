package plugins

import (
	"github.com/benishor/go-irc-bot/irc/bot"
	"regexp"
	"strings"
	"github.com/benishor/go-irc-bot/irc/commands"
	"fmt"
	"github.com/benishor/go-irc-bot/irc"
	"time"
)

var lastTool = NewLastSeenToolSql()

func init() {
	bot.RegisterPlugin(&LastSeenPlugin{})
}

var botCommandPattern = regexp.MustCompile(`!([^\W]+)(.*)`)

type LastSeenPlugin struct {
	bot.DefaultBotPlugin
}

func (p*LastSeenPlugin) OnJoinChannel(bot *bot.Bot, channel string) {
}

func (p*LastSeenPlugin) OnUserJoin(bot *bot.Bot, channel string, user string) {
}

func (p*LastSeenPlugin) OnUserPart(bot *bot.Bot, channel string, user string, partMessage string) {
	lastTool.set(user, channel, partMessage)
}

func (p*LastSeenPlugin) OnUserQuit(bot *bot.Bot, user string, quitMessage string) {
	lastTool.set(user, "", quitMessage)
}

func (p*LastSeenPlugin) OnChannelMessage(bot *bot.Bot, channel string, user string, message string) {
	result := botCommandPattern.FindStringSubmatch(message)
	if result == nil {
		return;
	}

	botCommand := result[1]
	if botCommand != "seen" {
		return;
	}

	issuer := irc.ParseIrcUser(user)

	searchedNickname := strings.TrimSpace(result[2])
	if searchedNickname == "" {
		bot.Write(commands.SendMessage(channel, fmt.Sprintf(`%s: correct command usage is "!seen [nickname]"`, issuer.Nickname)))
		return;
	} else if strings.ToLower(searchedNickname) == strings.ToLower(issuer.Nickname) {
		bot.Write(commands.SendMessage(channel, fmt.Sprintf(`%s: go look in a mirror!`, issuer.Nickname)))
		return;
	} else if strings.ToLower(searchedNickname) == strings.ToLower(bot.CurrentNickname) {
		bot.Write(commands.SendMessage(channel, fmt.Sprintf(`%s: I am right here!!1`, issuer.Nickname)))
		return;
	}

	// check if user on channel
	userCurrentlyOnChannel := false
	for _, user := range bot.GetChannelStats(channel).Users {
		if user.Nickname == searchedNickname {
			userCurrentlyOnChannel = true
			break
		}
	}
	if userCurrentlyOnChannel {
		bot.Write(commands.SendMessage(channel,
			fmt.Sprintf("%s: it seems like you need glasses. %s is in the channel!",
				issuer.Nickname, searchedNickname)))
		return;
	}

	lastSeenInfo := lastTool.get(searchedNickname)
	if lastSeenInfo == nil {
		bot.Write(commands.SendMessage(channel, fmt.Sprintf("%s: sorry, I didn't see %s", issuer.Nickname, searchedNickname)))
		return;
	}

	lastSeenMessage := ""
	if lastSeenInfo.Channel == "" {
		if lastSeenInfo.Message == "" {
			lastSeenMessage = fmt.Sprintf("%s: %s (%s) was seen quitting on %s",
				issuer.Nickname,
				searchedNickname,
				lastSeenInfo.User,
				lastSeenInfo.LeavingTime.Format(time.UnixDate))
		} else {
			lastSeenMessage = fmt.Sprintf("%s: %s (%s) was seen quitting on %s stating [%s]",
				issuer.Nickname,
				searchedNickname,
				lastSeenInfo.User,
				lastSeenInfo.LeavingTime.Format(time.UnixDate),
				lastSeenInfo.Message)
		}

	} else {
		if lastSeenInfo.Message == "" {
			lastSeenMessage = fmt.Sprintf("%s: %s (%s) was last seen leaving %s on %s",
				issuer.Nickname,
				searchedNickname,
				lastSeenInfo.User,
				lastSeenInfo.Channel,
				lastSeenInfo.LeavingTime.Format(time.UnixDate))
		} else {
			lastSeenMessage = fmt.Sprintf("%s: %s (%s) was last seen leaving %s on %s stating [%s]",
				issuer.Nickname,
				searchedNickname,
				lastSeenInfo.User,
				lastSeenInfo.Channel,
				lastSeenInfo.LeavingTime.Format(time.UnixDate),
				lastSeenInfo.Message)
		}

	}
	bot.Write(commands.SendMessage(channel, lastSeenMessage))
}
