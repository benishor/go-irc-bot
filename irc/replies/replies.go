package replies

const (
	CmdPing = "PING"
	CmdPrivMsg = "PRIVMSG"
)

const (
	RplWelcome = "001"
	RplYourHost = "002"
	RplCreated = "003"
	RplMyInfo = "004"
	RplBounce = "005"
	RplUserHost = "302"
	RplIsOn = "303"
	RplEndOfWhois = "318"
	RplNoTopic = "331"
	RplTopic = "332"
	RplNameReply = "353"
	RplEndOfNames = "366"
	RplMotdStart = "375"
	RplMotd = "372"
	RplEndOfMotd = "376"
	RplJoin = "JOIN"
	RplPart = "PART"
	RplQuit = "QUIT"
	RplMode = "MODE"
)

const (
	ErrNoMotd = "422"
	ErrErronoeusNickname = "432"
	ErrNicknameInUse = "433"
	ErrChannelIsFull = "471"
	ErrBannedFromChan = "474"
	ErrBadChannelKey = "475"
)
