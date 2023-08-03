package controllers

import (
	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"strconv"
)

func CreateUser(c *gin.Context) {
	var userInput models.UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username: userInput.Username,
		Password: userInput.Password,
		Email:    userInput.Email,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := models.UserResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}

	c.JSON(201, gin.H{"data": response})
}

func UpdateBalanceAndCredit(c *gin.Context) {
	var input struct {
		UserId  int `json:"user_id"`
		Balance int `json:"balance"`
		Credit  int `json:"credit"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	var detail models.Detail
	if err := tx.Where("id = ?", input.UserId).First(&detail).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if input.Credit > detail.Credit {
		tx.Rollback()
		c.JSON(400, gin.H{"error": "Credit amount exceeds the allowed credit limit"})
		return
	}

	if err := tx.Model(&models.Detail{}).Where("id = ?", input.UserId).Updates(models.Detail{Balance: input.Balance, Credit: input.Credit}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Balance and credit updated successfully"})
}

func GetAllUser(c *gin.Context) {
	var (
		userList []models.User
		pageSize = 12
		page     = 1
	)

	pageParam := c.Query("page")
	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	queryParams := c.Request.URL.Query()
	role := queryParams.Get("role")
	username := queryParams.Get("username")

	db := config.DB

	if role != "" && username != "" {
		modifiedUsername := "%" + username + "%"
		db = db.Where("username LIKE ?", modifiedUsername).Where("role_id = ?", role)
	} else if role != "" {
		db = db.Where("role_id = ?", role)
	} else if username != "" {
		modifiedUsername := "%" + username + "%"
		db = db.Where("username LIKE ?", modifiedUsername)
	}

	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	offset := (page - 1) * pageSize
	err := db.Model(&userList).
		Select("id", "username", "email", "created_at").
		Offset(offset).
		Limit(pageSize).
		Find(&userList).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	maxPage := (int(count) + pageSize - 1) / pageSize

	var userListResponse []models.UserResponse
	err = dto.Map(&userListResponse, &userList)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data":     userListResponse,
		"page":     page,
		"pageSize": pageSize,
		"total":    maxPage,
	})
}

func GetUserById(c *gin.Context) {
	var user models.User
	var userResponse models.UserResponse
	id := c.Param("id")

	err := config.DB.Where("users.id = ?", id).
		Select("username", "email", "detail_id").
		First(&user).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = dto.Map(&userResponse, &user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": userResponse,
	})
}

func GetAllUserWithoutPagination(c *gin.Context) {
	var userList []models.User
	err := config.DB.Where("role_id = 4").
		Joins("Role").
		Find(&userList).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"data": userList,
	})
}
