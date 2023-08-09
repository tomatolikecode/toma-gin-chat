package main

import (
	"github.com/toma-gin-chat/router"
	"github.com/toma-gin-chat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()

	router := router.Router()
	router.Run(":8080")
}
