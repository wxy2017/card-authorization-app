package handlers

import (
	"card-authorization/database"
	"card-authorization/log"
	"card-authorization/middleware"
	"card-authorization/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token   string       `json:"token"`
	User    *models.User `json:"user"`
	Message string       `json:"message"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "邮箱已存在"})
		return
	}

	// 创建新用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	}

	if err := database.DB.Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token:   token,
		User:    user,
		Message: "注册成功",
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:   token,
		User:    &user,
		Message: "登录成功",
	})
}

func GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func ListUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Table("users").Order("created_at desc").Find(&users).Error; err != nil {
		log.Error("获取用户失败", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "获取用户失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
func UpdateUser(c *gin.Context) {
	// 定义接收前端数据的结构体
	type Request struct {
		User models.User `json:"user"`
	}

	var req Request
	// 绑定并验证前端传递的JSON数据
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("数据绑定失败", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 检查用户ID是否存在
	if req.User.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不能为空"})
		return
	}

	// 执行更新操作
	if err := database.DB.Model(&models.User{}).Where("id = ?", req.User.ID).Updates(&req.User).Error; err != nil {
		log.Error("更新用户失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	// 返回更新成功的响应
	c.JSON(http.StatusOK, gin.H{
		"message": "用户信息更新成功",
		"user":    req.User,
	})
}

// LastActive 我最近的活动
func LastActive(c *gin.Context) {
	userID := c.GetUint("userID")
	userID = 1
	//最新发送给我的卡
	//最近我我创建的卡
	//最近我使用的卡
	//最近我发送给别人的卡

	var cards []models.Card
	if err := database.DB.Table("cards").
		Preload("Creator").
		Preload("Owner").
		Joins("INNER JOIN card_transactions ct ON ct.card_id = cards.id").
		Where("cards.status != ?", "expired").
		Where("ct.from_user_id = ? OR ct.to_user_id = ?", userID, userID).
		Select("cards.*, ct.created_at as transaction_at, ct.type as transaction_type").
		Order("ct.created_at DESC").
		Limit(5).
		Find(&cards).Error; err != nil {
		log.Error("数据查询异常: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据查询异常"})
		return
	}

	//格式化数据
	type Active struct {
		CreatorNickname string    `json:"creator_nickname"`
		OwnerNickname   string    `json:"owner_nickname"`
		CreatorID       uint      `json:"creator_id"`
		CardTitle       string    `json:"card_title"`
		CardDescription string    `json:"card_description"`
		TransactionAt   time.Time `json:"transaction_at"`
		TransactionType string    `json:"transaction_type"`
	}

	activeCards := make([]Active, 0, 5)
	for i := range cards {
		activeCards = append(activeCards, Active{
			CreatorNickname: cards[i].Creator.Nickname,
			OwnerNickname:   cards[i].Owner.Nickname,
			CreatorID:       cards[i].CreatorID,
			CardTitle:       cards[i].Title,
			CardDescription: cards[i].Description,
			TransactionAt:   *cards[i].TransactionAt,
			TransactionType: cards[i].TransactionType,
		})
	}

	// 返回更新成功的响应
	c.JSON(http.StatusOK, gin.H{
		"message":     "最近活动信息",
		"activeCards": activeCards,
	})
}

func Test(c *gin.Context) {
	LastActive(c)
}
