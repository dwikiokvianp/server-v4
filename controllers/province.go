package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func FindProvince(c *gin.Context) {
	province := models.Province{}

	if err := config.DB.Find(&province).Error; err != nil {
		c.JSON(404, gin.H{"message": "Not found"})
		return
	}

	c.JSON(200, province)

}

func FindProvinceById(c *gin.Context) {
	province := models.Province{}

	if err := config.DB.Where("id = ?", c.Param("id")).First(&province).Error; err != nil {
		c.JSON(404, gin.H{"message": "Not found"})
		return
	}

	c.JSON(200, province)

}
