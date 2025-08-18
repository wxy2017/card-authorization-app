package models

import (
	"time"
)

type CardStatus string

const (
	CardStatusActive  CardStatus = "active"
	CardStatusUsed    CardStatus = "used"
	CardStatusExpired CardStatus = "expired"
)

type Card struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Title           string     `gorm:"not null" json:"title"`
	Description     string     `gorm:"not null" json:"description"`
	CreatorID       uint       `gorm:"not null" json:"creator_id"`
	OwnerID         uint       `gorm:"not null" json:"owner_id"`
	Status          CardStatus `gorm:"default:active" json:"status"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	TransactionAt   *time.Time `json:"transaction_at,omitempty"`   // 交易时间
	TransactionType string     `json:"transaction_type,omitempty"` // 交易类型

	// 关联
	Creator User `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Owner   User `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

type CardTransaction struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CardID     uint      `gorm:"not null" json:"card_id"`
	FromUserID uint      `gorm:"not null" json:"from_user_id"`
	ToUserID   uint      `gorm:"not null" json:"to_user_id"`
	Type       string    `gorm:"not null" json:"type"` // send, use
	CreatedAt  time.Time `json:"created_at"`

	// 关联
	Card     Card `gorm:"foreignKey:CardID" json:"card,omitempty"`
	FromUser User `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUser   User `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
}
