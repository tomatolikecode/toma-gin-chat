package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/toma-gin-chat/docs"
	"github.com/toma-gin-chat/middleware"
	"github.com/toma-gin-chat/service"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.JWTAuth())
	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 静态资源
	// router.Static("/asset", "asset/")
	// router.LoadHTMLGlob("view/**/*")
	// router.GET("/", handlers ...gin.HandlerFunc)
	public := router.Group("")
	{
		public.GET("/index", service.Index)
		public.POST("/user/login", service.UserLogin)
		public.POST("/user/createUser", service.CreateUser)

		// 发送消息
		public.GET("/sendMsg", service.SendMsg)
		public.GET("/user/sendMsg", service.SendUserMsg)

	}
	private := router.Group("")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/user/getUserList", service.GetUserList)
		private.POST("/user/deleteUser", service.DeleteUser)
		private.POST("/user/updateUser", service.UpdateUser)
	}

	return router
}
