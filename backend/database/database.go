package database

import (
	"card-authorization/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("card_authorization.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移表结构
	err = DB.AutoMigrate(
		&models.User{},
		&models.Friends{},
		&models.Card{},
		&models.CardTransaction{},
	)
	if err != nil {
		return err
	}

	return nil
}
