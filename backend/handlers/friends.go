package handlers

import (
	"card-authorization/database"
	"card-authorization/log"
	"card-authorization/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListFriends 获取道友列表，按最近互动时间排序
func ListFriends(c *gin.Context) {
	userID := c.GetUint("userID")
	var users []models.User
	subQuery := database.DB.Table("friends").
		Select("friend_id").
		Where("user_id = ?", userID).
		Order("updated_at desc")

	if err := database.DB.Table("users").
		Select("id, username, nickname, email").
		Where("id IN (?)", subQuery).
		Find(&users).Error; err != nil {
		log.Error("获取道友失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取道友失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// ListMyInviteFriends 获取我邀请的好友列表
func ListMyInviteFriends(c *gin.Context) {
	userID := c.GetUint("userID")
	var users []models.User
	subQuery := database.DB.Table("friend_invites").
		Select("to_user_id").
		Where("from_user_id = ?", userID).
		Order("updated_at desc")

	if err := database.DB.Table("users").
		Select("id, username, nickname, email").
		Where("id IN (?)", subQuery).
		Find(&users).Error; err != nil {
		log.Error("获取好友邀请失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友邀请失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// ListInviteMyFriends 获取邀请我的道友列表
func ListInviteMyFriends(c *gin.Context) {
	userID := c.GetUint("userID")
	var users []models.User
	subQuery := database.DB.Table("friend_invites").
		Select("from_user_id").
		Where("to_user_id = ?", userID).
		Order("updated_at desc")

	if err := database.DB.Table("users").
		Select("id, username, nickname, email").
		Where("id IN (?)", subQuery).
		Find(&users).Error; err != nil {
		log.Error("获取好友邀请失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友邀请失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
