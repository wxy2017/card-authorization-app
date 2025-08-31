package handlers

import (
	"card-authorization/database"
	"card-authorization/log"
	"card-authorization/models"
	"net/http"
	"strings"

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

// SearchFriendUsers 搜索用户，并关联好友关系
func SearchFriendUsers(c *gin.Context) {
	query := c.Query("q")
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词至少2个字符"})
		return
	}
	//去除query首尾空格
	query = strings.TrimSpace(query)

	var users []models.User
	if err := database.DB.
		Where("username LIKE ? OR nickname LIKE ? OR email = ?", "%"+query+"%", "%"+query+"%", query).
		Select("id, username, nickname, email").
		Limit(25).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索用户失败"})
		return
	}
	//获取查询到的用户在朋友邀请表中是否有记录
	var friendInvites []models.FriendInvite
	// 需要把用户切片转换为 ID 切片，否则 GORM 无法展开结构体
	userIDs := make([]uint, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	if len(userIDs) > 0 { // 只有在有用户时才查询，避免 IN () 语法问题
		if err := database.DB.
			Where("from_user_id IN ?", userIDs).
			Find(&friendInvites).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友邀请关系失败"})
			return
		}
	}
	// 封装返回结构体，包含用户信息和邀请状态
	var results []gin.H
	for _, user := range users {
		invited := "default"
		for _, invite := range friendInvites {
			if invite.FromUserID == user.ID {
				// 如果invite.Status为空，给个默认值：default
				if invite.Status != "" {
					invited = invite.Status
				}
				break
			}
		}
		results = append(results, gin.H{
			"user":    user,
			"invited": invited,
		})
	}
	c.JSON(http.StatusOK, gin.H{"list": results})
}
