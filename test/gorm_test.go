package aaa

import (
	"testing"

	"github.com/toma-gin-chat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGorm(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:fanqie@tcp(43.139.116.74:3306)/toma-gin-chat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	d, err := db.DB()
	if err != nil {
		panic("failed to connect database")
	}
	if err := d.Ping(); err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(
		&models.UserBasic{},
		&models.Message{},
	)
}
