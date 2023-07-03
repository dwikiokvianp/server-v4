package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetDrivers(c *gin.Context) {
	var drivers []models.Driver

	if err := config.DB.Find(&drivers).Error; err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": drivers,
	})
}

func GetTransactionByDriverId(c *gin.Context) {
	id := c.Param("id")
	var transactions []models.Transaction

	if err := config.DB.Where("driver_id = ?", id).
		Joins("User.Company").
		Joins("Officer").
		Joins("Province").
		Joins("City").
		Find(&transactions).Error; err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": transactions,
	})
}
