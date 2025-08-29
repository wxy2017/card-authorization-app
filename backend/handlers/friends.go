package handlers

import (
	"card-authorization/database"
	"card-authorization/log"
	"card-authorization/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListFriends 获取道友列表，按最近互动时间排序
func ListFriends(c *gin.Context) {

	var users []models.User
	if err := database.DB.Table("users").
		Joins("inner joins friends on friends.user_id = users.id").
		Order("friends.updated_at desc").
		Find(&users).Error; err != nil {
		log.Error("获取道友失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取道友失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
