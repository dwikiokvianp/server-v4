package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func FindProvince(c *gin.Context) {
	var province []models.Province

	if err := config.DB.Preload("City").Find(&province).Error; err != nil {
		c.JSON(404, gin.H{"message": "Not found"})
		return
	}

	c.JSON(200, gin.H{
		"data": province,
	})

}

func FindProvinceById(c *gin.Context) {
	var province models.Province

	if err := config.DB.Preload("City").Where("id = ?", c.Param("id")).First(&province).Error; err != nil {
		c.JSON(404, gin.H{"message": "Not found"})
		return
	}

	c.JSON(200, gin.H{
		"data": province,
	})

}
