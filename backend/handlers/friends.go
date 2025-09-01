package handlers

import (
	"card-authorization/database"
	"card-authorization/log"
	"card-authorization/models"
	"card-authorization/utils"
	"net/http"
	"strings"
	"time"

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

	//获取查询到的用户在朋友邀请表中是否有记录
	var friendInvites []models.FriendInvite
	// 需要把用户切片转换为 ID 切片，否则 GORM 无法展开结构体
	userIDs := make([]uint, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	if len(userIDs) > 0 { // 只有在有用户时才查询，避免 IN () 语法问题
		if err := database.DB.
			Where("from_user_id = ? AND to_user_id IN ? ", userID, userIDs).
			Find(&friendInvites).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友邀请关系失败"})
			return
		}
	}
	// 封装返回结构体，包含用户信息和邀请状态
	var results []gin.H
	for _, user := range users {
		invited := "default"
		updatedAt := time.Now()
		for _, invite := range friendInvites {
			if invite.ToUserID == user.ID {
				// 如果invite.Status为空，给个默认值：default
				if invite.Status != "" {
					invited = invite.Status
				}
				updatedAt = invite.UpdatedAt
				break
			}
		}
		results = append(results, gin.H{
			"user":       user,
			"invited":    invited,
			"updated_at": updatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"list": results})
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

	//获取查询到的用户在朋友邀请表中是否有记录
	var friendInvites []models.FriendInvite
	// 需要把用户切片转换为 ID 切片，否则 GORM 无法展开结构体
	userIDs := make([]uint, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	if len(userIDs) > 0 { // 只有在有用户时才查询，避免 IN () 语法问题
		if err := database.DB.
			Where("to_user_id = ? AND from_user_id IN ?", userID, userIDs).
			Find(&friendInvites).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友邀请关系失败"})
			return
		}
	}
	// 封装返回结构体，包含用户信息和邀请状态
	var results []gin.H
	for _, user := range users {
		invited := "default"
		updatedAt := time.Now()
		for _, invite := range friendInvites {
			if invite.FromUserID == user.ID {
				// 如果invite.Status为空，给个默认值：default
				if invite.Status != "" {
					invited = invite.Status
				}
				updatedAt = invite.UpdatedAt
				break
			}
		}
		results = append(results, gin.H{
			"user":       user,
			"invited":    invited,
			"updated_at": updatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"list": results})
}

// ListFriendUsers 获取好友列表，按最近互动时间排序
func ListFriendUsers(c *gin.Context) {
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
		log.Error("获取好友失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})

}

// SearchFriendUsers 搜索用户，并关联好友关系
func SearchFriendUsers(c *gin.Context) {
	userID := c.GetUint("userID")
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
			Where("from_user_id = ? AND to_user_id IN ?", userID, userIDs).
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
			if invite.ToUserID == user.ID {
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

// InviteFriends 邀请好友
func InviteFriends(c *gin.Context) {
	//从URL中获取被邀请的用户ID
	inviteeID := c.Param("id")
	if inviteeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "被邀请用户ID不能为空"})
		return
	}
	var invitee models.User
	if err := database.DB.First(&invitee, inviteeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "被邀请用户不存在"})
		return
	}
	//获取当前用户ID
	userID := c.GetUint("userID")
	if userID == invitee.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能邀请自己为好友"})
		return
	}
	//检查是否已经是好友关系
	var existingFriend models.Friends
	if err := database.DB.Where("user_id = ? AND friend_id = ?", userID, invitee.ID).First(&existingFriend).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "你们已经是好友关系"})
		return
	}
	//检查是否已经发送过邀请
	var existingInvite models.FriendInvite
	if err := database.DB.Where("from_user_id = ? AND to_user_id = ?", userID, invitee.ID).First(&existingInvite).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "你已经邀请过该用户为好友"})
		return
	}
	//创建好友邀请记录
	invite := models.FriendInvite{
		FromUserID: userID,
		ToUserID:   invitee.ID,
		Status:     "pending", // 初始状态为等待
	}
	if err := database.DB.Create(&invite).Error; err != nil {
		log.Error("创建好友邀请失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建好友邀请失败"})
		return
	}
	//邮件通知被邀请用户
	var myUser models.User
	if err := database.DB.First(&myUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "当前用户不存在"})
		return
	}
	subject := "你有一个新的好友邀请"
	body := buildEmailBodyOfInviteFriend(myUser.Nickname, myUser.Email)
	if err := utils.SendEmail(invitee.Email, subject, body); err != nil {
		log.Error("发送好友邀请邮件失败", err)
		//返回成功响应
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "好友邀请已发送但邮件通知失败"})
	} else {
		//返回成功响应
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "好友邀请已发送"})
	}
}

// AcceptFriends 接受好友邀请
func AcceptFriends(c *gin.Context) {
	//从URL中获取邀请的用户ID
	inviterID := c.Param("id")
	if inviterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邀请用户ID不能为空"})
		return
	}
	var inviter models.User
	if err := database.DB.First(&inviter, inviterID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "邀请用户不存在"})
		return
	}
	//获取当前用户ID
	userID := c.GetUint("userID")
	if userID == inviter.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能接受自己的邀请"})
		return
	}
	//检查是否已经是好友关系
	var existingFriend models.Friends
	if err := database.DB.Where("user_id = ? AND friend_id = ?", userID, inviter.ID).First(&existingFriend).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "你们已经是好友关系"})
		return
	}
	//检查是否有未处理的邀请
	var invite models.FriendInvite
	if err := database.DB.Where("from_user_id = ? AND to_user_id = ? AND status = ?", inviter.ID, userID, "pending").First(&invite).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有找到未处理的好友邀请"})
		return
	}
	//更新邀请状态为已接受
	invite.Status = "accepted"
	if err := database.DB.Save(&invite).Error; err != nil {
		log.Error("更新好友邀请状态失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新好友邀请状态失败"})
		return
	}
	//创建双向好友关系记录
	friend1 := models.Friends{
		UserID:    userID,
		FriendID:  inviter.ID,
		UpdatedAt: time.Now(),
	}
	friend2 := models.Friends{
		UserID:    inviter.ID,
		FriendID:  userID,
		UpdatedAt: time.Now(),
	}
	if err := database.DB.Create(&friend1).Error; err != nil {
		log.Error("创建好友关系失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建好友关系失败"})
		return
	}
	if err := database.DB.Create(&friend2).Error; err != nil {
		log.Error("创建好友关系失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建好友关系失败"})
		return
	}
	//邮件通知被邀请用户
	var myUser models.User
	if err := database.DB.First(&myUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "当前用户不存在"})
		return
	}
	subject := myUser.Nickname + " 已接受你的道友邀请"
	body := buildEmailBodyOfAcceptFriend(myUser.Nickname, myUser.Email)
	if err := utils.SendEmail(inviter.Email, subject, body); err != nil {
		log.Error("接受道友邀请邮件发送失败", err)
		//返回成功响应
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "接受道友邀请成功，但邮件通知失败"})
	} else {
		//返回成功响应
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "好友邀请已接受"})
	}
}

// 生成美化的邮件内容（请求好友）
func buildEmailBodyOfInviteFriend(fromNickname, fromEmail string) string {
	body := `
        <!DOCTYPE html>
            <html>
            <head>
                <meta charset="UTF-8">
                <title>邀请通知</title>
                <style>
                    body {
                        font-family: 'Helvetica Neue', Arial, sans-serif;
                        background-color: #f9f9f9;
                        margin: 0;
                        padding: 20px;
                        color: #333;
                    }
                    .container {
                        max-width: 600px;
                        margin: 0 auto;
                        background-color: white;
                        border-radius: 12px;
                        box-shadow: 0 4px 12px rgba(0,0,0,0.1);
                        overflow: hidden;
                    }
                    .header {
                        background: linear-gradient(135deg, #4a90e2, #5c6bc0);
                        color: white;
                        padding: 25px 30px;
                        text-align: center;
                    }
                    .header h1 {
                        margin: 0;
                        font-size: 24px;
                        font-weight: 600;
                    }
                    .content {
                        padding: 30px;
                        text-align: center;
                    }
                    .greeting {
                        font-size: 18px;
                        margin-bottom: 25px;
                        color: #555;
                    }
                    .card-notification {
                        background-color: #fff8e1;
                        border-left: 5px solid #ffc107;
                        padding: 20px;
                        border-radius: 8px;
                        margin: 20px 0;
                        font-size: 16px;
                        line-height: 1.6;
                    }
                    .highlight {
                        color: #e91e63;
                        font-weight: bold;
                        font-size: 18px;
                    }
                    .app-link {
                        margin: 30px 0;
                        padding: 20px;
                        background-color: #e3f2fd;
                        border-radius: 8px;
                    }
                    .app-link a {
                        color: #1976d2;
                        font-size: 18px;
                        font-weight: bold;
                        text-decoration: none;
                        border-bottom: 2px solid #1976d2;
                        padding-bottom: 3px;
                    }
                    .app-link a:hover {
                        color: #0d47a1;
                        border-bottom-color: #0d47a1;
                    }
                    .footer {
                        background-color: #f5f5f5;
                        padding: 20px 30px;
                        text-align: center;
                        color: #777;
                        font-size: 14px;
                    }
                </style>
            </head>
            <body>
                <div class="container">
                    <div class="header">
                        <h1>🎉 邀请通知</h1>
                    </div>
                    <div class="content">
                        <p class="greeting">你好！</p>
                        <div class="card-notification">
							<span class="highlight">` + fromNickname + `</span>向你发送了道友申请：
                            <br><br>
                            <span class="highlight">` + fromEmail + `</span>
                        </div>
                        <div class="app-link">
                            点击访问应用查看详情：<br><br>
                            <a href="http://wangxiang-pro.top:18080/" target="_blank">点我查看吆🎀</a>
                        </div>
                        <p>快去体验专为情侣和朋友设计的互动卡片系统吧～</p>
                    </div>
                    <div class="footer">
                        这是一封自动发送的通知邮件，无需回复
                    </div>
                </div>
            </body>
        </html>
    `
	return body
}

// 生成美化的邮件内容（接受好友的邀请）
func buildEmailBodyOfAcceptFriend(fromNickname, fromEmail string) string {
	body := `
        <!DOCTYPE html>
            <html>
            <head>
                <meta charset="UTF-8">
                <title>邀请通过</title>
                <style>
                    body {
                        font-family: 'Helvetica Neue', Arial, sans-serif;
                        background-color: #f9f9f9;
                        margin: 0;
                        padding: 20px;
                        color: #333;
                    }
                    .container {
                        max-width: 600px;
                        margin: 0 auto;
                        background-color: white;
                        border-radius: 12px;
                        box-shadow: 0 4px 12px rgba(0,0,0,0.1);
                        overflow: hidden;
                    }
                    .header {
                        background: linear-gradient(135deg, #4a90e2, #5c6bc0);
                        color: white;
                        padding: 25px 30px;
                        text-align: center;
                    }
                    .header h1 {
                        margin: 0;
                        font-size: 24px;
                        font-weight: 600;
                    }
                    .content {
                        padding: 30px;
                        text-align: center;
                    }
                    .greeting {
                        font-size: 18px;
                        margin-bottom: 25px;
                        color: #555;
                    }
                    .card-notification {
                        background-color: #fff8e1;
                        border-left: 5px solid #ffc107;
                        padding: 20px;
                        border-radius: 8px;
                        margin: 20px 0;
                        font-size: 16px;
                        line-height: 1.6;
                    }
                    .highlight {
                        color: #e91e63;
                        font-weight: bold;
                        font-size: 18px;
                    }
                    .app-link {
                        margin: 30px 0;
                        padding: 20px;
                        background-color: #e3f2fd;
                        border-radius: 8px;
                    }
                    .app-link a {
                        color: #1976d2;
                        font-size: 18px;
                        font-weight: bold;
                        text-decoration: none;
                        border-bottom: 2px solid #1976d2;
                        padding-bottom: 3px;
                    }
                    .app-link a:hover {
                        color: #0d47a1;
                        border-bottom-color: #0d47a1;
                    }
                    .footer {
                        background-color: #f5f5f5;
                        padding: 20px 30px;
                        text-align: center;
                        color: #777;
                        font-size: 14px;
                    }
                </style>
            </head>
            <body>
                <div class="container">
                    <div class="header">
                        <h1>🎉 邀请通过</h1>
                    </div>
                    <div class="content">
                        <p class="greeting">你好！</p>
                        <div class="card-notification">
							<span class="highlight">` + fromNickname + `</span>同意了你的道友申请
                            <br><br>
                            <span class="highlight">` + fromEmail + `</span>
                        </div>
                        <div class="app-link">
                            点击访问应用查看详情：<br><br>
                            <a href="http://wangxiang-pro.top:18080/" target="_blank">点我查看吆🎀</a>
                        </div>
                        <p>快去体验专为情侣和朋友设计的互动卡片系统吧～</p>
                    </div>
                    <div class="footer">
                        这是一封自动发送的通知邮件，无需回复
                    </div>
                </div>
            </body>
        </html>
    `
	return body
}
