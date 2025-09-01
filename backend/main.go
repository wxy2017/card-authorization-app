package main

import (
	"card-authorization/config"
	"card-authorization/database"
	"card-authorization/handlers"
	"card-authorization/log"
	"card-authorization/middleware"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}
	//加载外部配置文件
	if err := config.LoadConfig(); err != nil {
		log.Fatal("加载外部配置文件失败:", err)
	}

	// 创建Gin路由
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	// 移除默认日志中间件
	r.Use(gin.Recovery())
	// 注册自定义日志中间件
	r.Use(customGinLogger(strconv.Itoa(os.Getpid())))

	// 静态文件服务
	r.Static("/static", "../frontend/static")
	r.LoadHTMLGlob("../frontend/templates/*")

	// 前端页面路由
	r.GET("/", handlers.Index)
	r.GET("/login", handlers.LoginPage)
	r.GET("/register", handlers.RegisterPage)
	r.GET("/dashboard", handlers.Dashboard)
	r.GET("/cards", handlers.CardsPage)
	r.GET("/cards/create", handlers.CreateCardPage)
	r.GET("/friends", handlers.Friends)

	// API路由组
	api := r.Group("/api")
	{
		// 用户相关(无需鉴权)
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.GET("/test", handlers.Test)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.AuthRequired())
		{
			// 卡片相关
			auth.POST("/cards", handlers.CreateCard)
			auth.GET("/cards", handlers.GetMyCards)
			auth.GET("/cards/received", handlers.GetReceivedCards)
			auth.GET("/cards/send", handlers.GetSendCards)
			auth.POST("/cards/used", handlers.UsedCard)
			auth.POST("/cards/:id/use", handlers.UseCard)
			auth.POST("/cards/:id/send", handlers.SendCard)
			auth.POST("/cards/:id/delete", handlers.DeleteCard)
			auth.GET("/cards/:id/copy", handlers.CopyCard)
			// 用户相关
			auth.GET("/profile", handlers.GetProfile)
			auth.GET("/users/listUsers", handlers.ListUsers)
			auth.POST("/user/:id/update", handlers.UpdateUser)
			auth.GET("/users/lastActive", handlers.LastActive)
			auth.GET("/users/friends", handlers.ListFriends)
			auth.GET("/users/friends/search", handlers.SearchFriendUsers)
			auth.GET("/users/friends/list", handlers.ListFriendUsers)
			auth.GET("/users/friends/myInvite/list", handlers.ListMyInviteFriends)
			auth.GET("/users/friends/inviteMy/list", handlers.ListInviteMyFriends)
			auth.GET("/users/friends/:id/invite", handlers.InviteFriend)
			auth.GET("/users/friends/:id/accept", handlers.AcceptFriend)
			auth.GET("/users/friends/:id/delete", handlers.DeleteFriend)
		}
	}

	//启动定时器
	go handlers.CheckExpiredCards()

	// 启动服务器
	log.Info("服务器启动在 http://localhost:" + config.SystemConfig.HTTPPort)
	if err := r.Run(":" + config.SystemConfig.HTTPPort); err != nil {
		log.Fatal("服务器启动失败:", err)
	}

}

// customGinLogger 创建自定义日志中间件
func customGinLogger(pid string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 计算耗时
		latency := time.Since(startTime)
		// 获取请求信息
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		// 格式化日志内容
		logMsg := fmt.Sprintf("[%s]: [GIN] %s - %s | %d | %12s | %15s | %-6s \"%s\"",
			pid,
			startTime.Format("2006/01/02"),
			startTime.Format("15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
		// 使用自定义的log.Info输出
		log.Info(logMsg)
	}
}
