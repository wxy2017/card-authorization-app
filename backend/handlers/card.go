package handlers

import (
	"card-authorization/database"
	"card-authorization/log"
	"card-authorization/models"
	"card-authorization/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateCardRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"required"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type SendCardRequest struct {
	ToUsername string `json:"to_username" binding:"required"`
}

func CreateCard(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   userID,
		OwnerID:     userID,
		ExpiresAt:   req.ExpiresAt,
	}

	if err := database.DB.Create(card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建卡片失败"})
		return
	}

	// 预加载关联数据
	database.DB.Preload("Creator").First(card, card.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "卡片创建成功",
		"card":    card,
	})
}

func GetMyCards(c *gin.Context) {
	userID := c.GetUint("userID")

	var cards []models.Card
	if err := database.DB.Preload("Creator").Preload("Owner").
		Where("creator_id = ?", userID).
		Order("created_at DESC").
		Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取卡片失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cards": cards})
}

func GetReceivedCards(c *gin.Context) {
	userID := c.GetUint("userID")

	var cards []models.Card
	if err := database.DB.Preload("Creator").Preload("Owner").
		Where("owner_id = ? AND creator_id != ?", userID, userID).
		Order("created_at DESC").
		Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取卡片失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cards": cards})
}

func UseCard(c *gin.Context) {
	userID := c.GetUint("userID")
	cardID := c.Param("id")

	var card models.Card
	if err := database.DB.Preload("Creator").Preload("Owner").First(&card, cardID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "卡片不存在"})
		return
	}

	// 检查卡片所有者
	if card.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权使用该卡片"})
		return
	}

	// 检查卡片状态
	if card.Status != models.CardStatusActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "卡片已使用或已过期"})
		return
	}

	// 更新卡片状态
	card.Status = models.CardStatusUsed
	if err := database.DB.Save(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新卡片状态失败"})
		return
	}

	// 记录交易
	transaction := &models.CardTransaction{
		CardID:     card.ID,
		FromUserID: card.OwnerID,
		ToUserID:   card.CreatorID,
		Type:       "use",
	}
	if err := database.DB.Create(transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录交易失败"})
		return
	}
	//发送邮件通知卡片创造者，拥有者已经使用当前卡片 todo:这里有bug还没有写好
	if card.Creator.Email != "" {
		var body = buildEmailBody(card.Owner.Nickname, card.Title)
		if err := utils.SendEmail(card.Creator.Email, "卡片已被使用："+card.Title, body); err != nil {
			log.Error("向%s发送邮件失败", card.Creator.Nickname)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "卡片使用成功",
		"card":    card,
	})
}

func SendCard(c *gin.Context) {
	userID := c.GetUint("userID")
	cardID := c.Param("id")

	var req SendCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var card models.Card
	if err := database.DB.Preload("Owner").First(&card, cardID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "卡片不存在"})
		return
	}

	// 检查卡片所有者
	if card.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "所属者非本人，无权发送该卡片"})
		return
	}

	// 检查卡片状态
	if card.Status != models.CardStatusActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "卡片已使用或已过期"})
		return
	}

	// 查找接收用户
	var toUser models.User
	if err := database.DB.Where("username = ?", req.ToUsername).First(&toUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "接收用户不存在"})
		return
	}

	// 更新卡片所有者
	oldOwner := card.Owner
	card.OwnerID = toUser.ID
	card.Owner = models.User{}
	if err := database.DB.Save(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送卡片失败"})
		return
	}

	// 记录交易
	transaction := &models.CardTransaction{
		CardID:     card.ID,
		FromUserID: userID,
		ToUserID:   toUser.ID,
		Type:       "send",
	}
	if err := database.DB.Create(transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录交易失败"})
		return
	}
	// 如果接收者有邮箱则发送邮件
	if toUser.Email != "" {
		// Email 不为空的逻辑
		//var body = "恭喜你，收到来自" + oldOwner.Nickname + "的卡：" + card.Title
		var body = buildEmailBody(oldOwner.Nickname, card.Title)
		if err := utils.SendEmail(toUser.Email, "收到卡："+card.Title, body); err != nil {
			log.Error("向%s发送邮件失败", toUser.Nickname)
			c.JSON(http.StatusOK, gin.H{
				"message": "卡片发送成功，邮件通知失败！",
				"card":    card,
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "卡片发送成功，已邮件通知！",
				"card":    card,
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "卡片发送成功",
			"card":    card,
		})
	}
}

func SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词至少2个字符"})
		return
	}

	var users []models.User
	if err := database.DB.
		Where("username LIKE ? OR nickname LIKE ?", "%"+query+"%", "%"+query+"%").
		Select("id, username, nickname").
		Limit(10).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func DeleteCard(c *gin.Context) {
	userID := c.GetUint("userID")
	cardID := c.Param("id")
	var card models.Card
	if err := database.DB.First(&card, cardID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "卡片不存在"})
		return
	}
	// 检查卡片所有者
	if card.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除该卡片"})
		return
	}
	// 检查卡片创造者
	if card.CreatorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除该卡片"})
		return
	}
	if err := database.DB.Delete(&models.Card{}, cardID).Error; err != nil {
		log.Error("卡%s删除失败", cardID)
		c.JSON(http.StatusExpectationFailed, gin.H{"error": "删除失败"})
		return
	}
	log.Error("卡[%d:%s]删除成功", card.ID, card.Title)
	c.JSON(http.StatusOK, gin.H{"message": "卡片删除成功"})
}

// 生成美化的邮件内容
func buildEmailBody(formNickname, cardTitle string) string {
	body := `
        <!DOCTYPE html>
            <html>
            <head>
                <meta charset="UTF-8">
                <title>新卡片通知</title>
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
                        <h1>🎉 新卡片通知</h1>
                    </div>
                    <div class="content">
                        <p class="greeting">恭喜你！</p>
                        <div class="card-notification">
                            你收到了来自 <span class="highlight">` + formNickname + `</span> 的卡：
                            <br><br>
                            <span class="highlight">` + cardTitle + `</span>
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
