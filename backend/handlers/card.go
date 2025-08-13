package handlers

import (
	"card-authorization/database"
	"card-authorization/log"
	"card-authorization/models"
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
	if err := database.DB.First(&card, cardID).Error; err != nil {
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
	if err := database.DB.First(&card, cardID).Error; err != nil {
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
	card.OwnerID = toUser.ID
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

	c.JSON(http.StatusOK, gin.H{
		"message": "卡片发送成功",
		"card":    card,
	})
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
