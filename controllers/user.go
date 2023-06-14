package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"strconv"
)

func GetAllUser(c *gin.Context) {
	var (
		userList []models.User
		pageSize = 10
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
	}
	if username == "" {
		db = db.Where("role_id = ?", role)
	}

	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	offset := (page - 1) * pageSize

	err := db.Offset(offset).Limit(pageSize).Preload("Role").Preload("Detail").Find(&userList).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	page = (int(count) + pageSize - 1) / pageSize

	c.JSON(200, gin.H{
		"data":     userList,
		"page":     page,
		"pageSize": pageSize,
		"total":    page,
	})
}

func GetUserById(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	err := config.DB.Where("id = ?", id).Preload("Role").Preload("Detail").Preload("Company").First(&user).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"data": user,
	})
}
