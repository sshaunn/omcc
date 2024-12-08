package common

import tele "gopkg.in/telebot.v3"

type SocialPlatformType int
type TradingPlatformType int
type Status string
type MemberStatus tele.MemberStatus

// General ENV constants
const (
	BitgetApiKeyEnvPath        string = "exchange.bitget.apiKey"
	BitgetApiSecretKeyEnvPath         = "exchange.bitget.secretKey"
	BitgetApiPassphraseEnvPath        = "exchange.bitget.passphrase"
)

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
	DailyTrading   string = "daily"
	WeeklyTrading         = "weekly"
	MonthlyTrading        = "monthly"
)

const (
	Normal      Status = "normal"
	Whitelisted        = "whitelisted"
	Blacklisted        = "blacklisted"
)

const (
	Creator       MemberStatus = "creator"
	Administrator              = "administrator"
	Member                     = "member"
	Restricted                 = "restricted"
	Left                       = "left"
	Kicked                     = "kicked"
	Unknown                    = "unknown"
)

var statusMap = map[string]MemberStatus{
	"creator":       Creator,
	"administrator": Administrator,
	"member":        Member,
	"restricted":    Restricted,
	"left":          Left,
	"kicked":        Kicked,
}

const (
	Private int = iota
	Supergroup
	Group
	Channel
)

func (m MemberStatus) Value() string {
	switch m {
	case Creator:
		return "擁有者"
	case Administrator:
		return "管理員"
	case Member:
		return "成員"
	case Restricted:
		return "有限"
	case Left:
		return "已退出"
	case Kicked:
		return "已封鎖"
	default:
		return "未知"
	}
}

func GetMemberStatusFromValue(value tele.MemberStatus) MemberStatus {
	if status, ok := statusMap[string(value)]; ok {
		return status
	}
	return Unknown
}

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
	MemberStatus   MemberStatus
	SocialPlatform SocialPlatformType
}
