package main

import (
	"DEMO01/controller"
	"DEMO01/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	handler.InitDB()
	handler.RedisInit()
	defer handler.DB.Close()
	defer handler.Rdb.Close()
	r := gin.Default()
	//r.Use(controller.Cors)
	r.Use(handler.JwtVerify)
	
	usegroup := r.Group("/")
	{
		usegroup.POST("login", controller.HandleLogin)
		usegroup.POST("register", controller.HandleRegister)
		usegroup.POST("register/captcha", controller.HandleRegisterCaptcha)
		usegroup.GET("userinfo", controller.HandleUserinfo)
	}
	contentgroup := r.Group("/content")
	{
		contentgroup.GET("/channels", controller.HandleChannels)
		contentgroup.GET("/articles",controller.HandleArticles)
		contentgroup.POST("/articles",controller.HandleUpload)
		contentgroup.PUT("/articles",controller.HandleUpdate)
		contentgroup.POST("/upload",controller.HandleImagesUpload)
	}
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

