package main

import (
	"card-authorization/database"
	"card-authorization/handlers"
	"card-authorization/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 创建Gin路由
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 静态文件服务
	r.Static("/static", "../frontend/static")
	r.LoadHTMLGlob("../frontend/templates/*")

	// API路由组
	api := r.Group("/api")
	{
		// 用户相关
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.AuthRequired())
		{
			// 卡片相关
			auth.POST("/cards", handlers.CreateCard)
			auth.GET("/cards", handlers.GetMyCards)
			auth.GET("/cards/received", handlers.GetReceivedCards)
			auth.POST("/cards/:id/use", handlers.UseCard)
			auth.POST("/cards/:id/send", handlers.SendCard)

			// 用户相关
			auth.GET("/profile", handlers.GetProfile)
			auth.GET("/users/search", handlers.SearchUsers)
		}
	}

	// 前端页面路由
	r.GET("/", handlers.Index)
	r.GET("/login", handlers.LoginPage)
	r.GET("/register", handlers.RegisterPage)
	r.GET("/dashboard", handlers.Dashboard)
	r.GET("/cards", handlers.CardsPage)
	r.GET("/cards/create", handlers.CreateCardPage)

	// 启动服务器
	log.Println("服务器启动在 http://localhost:18080")
	if err := r.Run(":18080"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
