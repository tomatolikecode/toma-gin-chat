package utils

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func InitMysql() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		})
	dns := viper.GetString("mysql.dns")
	if dns == "" {
		panic("db connect error, dsn is nil")
	}
	Db, _ = gorm.Open(mysql.Open(dns), &gorm.Config{Logger: newLogger})
	d, _ := Db.DB()
	if err := d.Ping(); err != nil {
		panic("db connect error, " + err.Error())
	}
}
