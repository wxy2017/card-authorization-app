package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Friends struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	FriendID  uint      `json:"friend_id"`
	UpdatedAt time.Time `json:"created_at"`
}

type FriendInvite struct {
	ID         uint      `json:"id"`
	FromUserID uint      `json:"from_user_id"`
	ToUserID   uint      `json:"to_user_id"`
	Status     string    `json:"status"` // pending（等待）, accepted（接收）, rejected（拒接）
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// 创建用户前加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
