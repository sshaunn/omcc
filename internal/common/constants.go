package common

type SocialPlatformType int
type TradingPlatformType int
type Status string

// Social platform constants
const (
	_ SocialPlatformType = iota
	Telegram
	Line
)

const (
	_ TradingPlatformType = iota
	Bitget
	BingX
)

const (
	Normal      Status = "normal"
	Whitelisted        = "whitelisted"
	Blacklisted        = "blacklisted"
)

func (p SocialPlatformType) Name() string {
	return [...]string{"TELEGRAM", "LINE"}[p]
}

func (p SocialPlatformType) Value() int {
	return int(p)
}

func (t TradingPlatformType) Name() string {
	return [...]string{"BITGET", "BINGX"}[t]
}

func (t TradingPlatformType) Value() int {
	return int(t)
}

// UserInfo 存储验证后的用户信息
type UserInfo struct {
	UID            string
	UserId         string
	Firstname      string
	Lastname       string
	Username       string
	SocialPlatform SocialPlatformType
}
