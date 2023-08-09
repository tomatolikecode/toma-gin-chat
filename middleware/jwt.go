package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/toma-gin-chat/utils"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			// 先跳过
			c.Next()
			return
			// c.JSON(401, gin.H{"code": 401, "msg": "登录失效"})
			// c.Abort()
			// return
		}
		j := utils.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			c.JSON(401, gin.H{"code": 401, "msg": "登录失效"})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
