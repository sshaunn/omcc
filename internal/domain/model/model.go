package model

import (
	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"time"
)

type Customer struct {
	Id        string    `gorm:"primary_key;type:varchar(36)" json:"id"`
	Username  string    `gorm:"type:varchar(50)" json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SocialPlatform struct {
	Id       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(50)" json:"name"`
	IsActive bool   `gorm:"type:bool;default:true" json:"is_active"`
}

type TradingPlatform struct {
	Id       int    `gorm:"primary_key;autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(50)" json:"name"`
	IsActive bool   `gorm:"type:bool;default:true" json:"is_active"`
}

type CustomerSocialBinding struct {
	ID            int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID    string            `gorm:"type:varchar(36)" json:"customer_id"`
	SocialID      int               `gorm:"type:int" json:"social_id"`
	UserID        string            `gorm:"type:varchar(50)" json:"user_id"`
	Username      string            `gorm:"type:varchar(50)" json:"username"`
	Firstname     string            `gorm:"type:varchar(50)" json:"firstname"`
	Lastname      string            `gorm:"type:varchar(50)" json:"lastname"`
	IsActive      bool              `gorm:"default:true" json:"is_active"`
	DeactivatedAt *time.Time        `json:"deactivated_at"`
	MemberStatus  tele.MemberStatus `gorm:"type:enum('creator', 'administrator', 'member', 'restricted', 'left', 'kicked')" json:"member_status"`
	Status        common.Status     `gorm:"type:enum('normal','whitelisted','blacklisted')" json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Customer      *Customer         `gorm:"foreignKey:CustomerID" json:"-"`
	Platform      *SocialPlatform   `gorm:"foreignKey:SocialID" json:"-"`
}

type CustomerTradingBinding struct {
	ID           int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID   string           `gorm:"type:varchar(36)" json:"customer_id"`
	TradingID    int              `json:"trading_id"`
	UID          string           `gorm:"type:varchar(50)" json:"uid"`
	RegisterTime string           `json:"register_time"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Customer     *Customer        `gorm:"foreignKey:CustomerID" json:"-"`
	Platform     *TradingPlatform `gorm:"foreignKey:TradingID" json:"-"`
}

type TradingHistory struct {
	ID          int64                  `gorm:"primaryKey;autoIncrement" json:"id"`
	BindingID   int64                  `json:"binding_id"`
	Volume      float64                `gorm:"type:decimal(16,2)" json:"volume"`
	TimePeriod  string                 `gorm:"type:enum('daily','weekly','monthly')" json:"time_period"`
	TradingDate time.Time              `json:"trading_date"`
	Binding     CustomerTradingBinding `gorm:"foreignKey:BindingID" json:"-"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	c.Id = uuid.New().String()
	return nil
}
