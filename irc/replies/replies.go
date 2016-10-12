package replies

const (
	CmdPing = "PING"
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
	RplMotdStart = "375"
	RplMotd = "372"
	RplEndOfMotd = "376"
	RplEndOfNames = "366"
	RplJoin = "JOIN"
)

const (
	ErrNoMotd = "422"
	ErrErronoeusNickname = "432"
	ErrNicknameInUse = "433"
)
