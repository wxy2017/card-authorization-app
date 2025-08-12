package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "功能卡片授权",
	})
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "登录 - 功能卡片授权",
	})
}

func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "注册 - 功能卡片授权",
	})
}

func Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "控制台 - 功能卡片授权",
	})
}

func CardsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "cards.html", gin.H{
		"title": "我的卡片 - 功能卡片授权",
	})
}

func CreateCardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_card.html", gin.H{
		"title": "创建卡片 - 功能卡片授权",
	})
}
