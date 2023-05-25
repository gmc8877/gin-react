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
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r.Use(controller.Cors)
	
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
		contentgroup.GET("/articles",controller.HandleUsrArticlesList)
		contentgroup.GET("/articles/:id",controller.HandleArticles)
		contentgroup.POST("/articles",controller.HandleUpload)
		contentgroup.PUT("/articles/:id",controller.HandleUpdate)
		contentgroup.POST("/upload",controller.HandleImagesUpload)
		contentgroup.DELETE("/articles/:id",controller.HandleUsrDelete)
	}
	r.GET("/channels", controller.HandleChannels)

	rootgroup := r.Group("/root")
	{
		rootgroup.POST("/login", controller.HandleRootLogin)
		rootgroup.Use(handler.RootVerify)
		rootgroup.GET("/articles",controller.HandleArticlesList)
		rootgroup.GET("/articles/:id",controller.HandleArticles)
		rootgroup.POST("/articles",controller.HandleUpload)
		rootgroup.PUT("/articles/:id",controller.HandleUpdate)
		rootgroup.POST("/upload",controller.HandleImagesUpload)
		rootgroup.DELETE("/articles/:id",controller.HandleDelete)
		rootgroup.GET("/rootinfo", controller.HandleRootinfo)
	}

	showgroup := r.Group("/home")
	{
		showgroup.GET("/articles", controller.HandleShows)
	}
	r.Run("localhost:8080") // 监听并在 0.0.0.0:8080 上启动服务
}

