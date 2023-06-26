package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func FindCity(c *gin.Context) {
	var city []models.City

	if err := config.DB.Find(&city).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Error",
			"error":   err,
		})
		return
	}

	c.JSON(200, gin.H{
		"data": city,
	})
}

func FindCityById(c *gin.Context) {
	id := c.Param("id")
	city := models.City{}

	if err := config.DB.Where("id = ?", id).Find(&city).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Error",
			"error":   err,
		})
		return
	}

	c.JSON(200, gin.H{
		"data": city,
	})
}
