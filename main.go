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
	r.Use(controller.Cors)
	
	r.GET("/assets/:path",  controller.Getimages)
	usegroup := r.Group("/")
	{
		usegroup.Use(handler.JwtVerify)
		usegroup.POST("login", controller.HandleLogin)
		usegroup.POST("register", controller.HandleRegister)
		usegroup.POST("register/captcha", controller.HandleRegisterCaptcha)
		usegroup.GET("userinfo", controller.HandleUserinfo)
	}
	contentgroup := r.Group("/content")
	{
		contentgroup.Use(handler.JwtVerify)
		contentgroup.GET("/channels", controller.HandleChannels)
		contentgroup.GET("/articles",controller.HandleArticlesList)
		contentgroup.GET("/articles/:id",controller.HandleArticles)
		contentgroup.POST("/articles",controller.HandleUpload)
		contentgroup.PUT("/articles/:id",controller.HandleUpdate)
		contentgroup.POST("/upload",controller.HandleImagesUpload)
		contentgroup.DELETE("/articles/:id",controller.HandleDelete)
	}
	showgroup := r.Group("/home")
	{
		showgroup.GET("/articles", controller.HandleShows)
	}
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

