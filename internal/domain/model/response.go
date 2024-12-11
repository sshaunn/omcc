package model

import (
	"math"
	"time"
)

type CustomerInfoResponse struct {
	Customer           CustomerInfo        `json:"customer"`
	SocialAccountInfo  CustomerSocialInfo  `json:"social_account_info"`
	TradingAccountInfo CustomerTradingInfo `json:"trading_account_info"`
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

type CustomerWithBindings struct {
	// Customer fields
	CustomerId        string    `gorm:"column:c_id"`
	CustomerUsername  string    `gorm:"column:c_username"`
	CustomerCreatedAt time.Time `gorm:"column:c_created_at"`
	CustomerUpdatedAt time.Time `gorm:"column:c_updated_at"`

	// Social binding fields
	SocialUserId       string    `gorm:"column:s_user_id"`
	SocialUsername     string    `gorm:"column:s_username"`
	SocialFirstname    string    `gorm:"column:s_firstname"`
	SocialLastname     string    `gorm:"column:s_lastname"`
	SocialIsActive     bool      `gorm:"column:s_is_active"`
	SocialStatus       string    `gorm:"column:s_status"`
	SocialMemberStatus string    `gorm:"column:s_member_status"`
	SocialCreatedAt    time.Time `gorm:"column:s_created_at"`

	// Trading binding fields
	TradingUid          string    `gorm:"column:t_uid"`
	TradingRegisterTime string    `gorm:"column:t_register_time"`
	TradingCreatedAt    time.Time `gorm:"column:t_created_at"`
}

type CustomerInfo struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type CustomerSocialInfo struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	IsActive     bool      `json:"is_active"`
	Status       string    `json:"status"`
	MemberStatus string    `json:"member_status"`
	SocialType   string    `json:"social_type"`
	CreatedAt    time.Time `json:"created_at"`
}

type CustomerTradingInfo struct {
	UID          string    `json:"uid"`
	RegisterTime string    `json:"register_time"`
	Platform     string    `json:"platform"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewPaginatedResponse[T any](data []T, total int64, page, limit int) *PaginatedResponse[T] {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &PaginatedResponse[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
