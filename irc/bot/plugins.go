package bot

// Interface that all bot plugins must implement
type BotPlugin interface {
	OnJoinChannel(bot *Bot, channel string)
	OnUserJoin(bot *Bot, channel string, user string)
	OnUserPart(bot *Bot, channel string, user string, partMessage string)
	OnUserQuit(bot *Bot, user string, quitMessage string)
	OnChannelMessage(bot *Bot, channel string, user string, message string)
}

// Holds the self-registering plugins
var Plugins []BotPlugin

func RegisterPlugin(pluginInstance BotPlugin) {
	Plugins = append(Plugins, pluginInstance)
}

// Empty plugin which does nothing but help other plugins "inherit" and override only what they need to
type DefaultBotPlugin struct{}

func (p*DefaultBotPlugin) OnJoinChannel(bot *Bot, channel string) {}
func (p*DefaultBotPlugin) OnUserJoin(bot *Bot, channel string, user string) {}
func (p*DefaultBotPlugin) OnUserPart(bot *Bot, channel string, user string, partMessage string) {}
func (p*DefaultBotPlugin) OnUserQuit(bot *Bot, user string, quitMessage string) {}
func (p*DefaultBotPlugin) OnChannelMessage(bot *Bot, channel string, user string, message string) {}
