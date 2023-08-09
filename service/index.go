package service

import (
	"github.com/gin-gonic/gin"
	"github.com/toma-gin-chat/utils"
)

// PingExample godoc
// @Summary index test
// @Description index hello
// @Tags index
// @Success 200 {string} string
// @Router /index [get]
func Index(c *gin.Context) {
	utils.Publish(c, utils.PublishKey, "测试啊")
	OkWithMsg("PONG", c)
}
