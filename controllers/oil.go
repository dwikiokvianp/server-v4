package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"time"
)

func GetAllOil(c *gin.Context) {
	var oil []models.Oil

	config.DB.Find(&oil)

	c.JSON(200, gin.H{
		"data": oil,
	})
}

func UpdateOilQuantity(c *gin.Context) {
	var inputOil models.Oil
	var oil models.Oil

	if err := c.ShouldBindJSON(&inputOil); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if inputOil.Quantity < 0 {
		c.JSON(400, gin.H{
			"error": "Quantity must be greater than 0",
		})
		return
	}

	config.DB.First(&oil, inputOil.Id)

	oil.Quantity = oil.Quantity + inputOil.Quantity

	config.DB.Save(&oil)

	var history models.HistoryIn
	history.OilId = oil.Id
	history.Quantity = inputOil.Quantity
	history.Date = time.Now()

	config.DB.Create(&history)

	c.JSON(200, gin.H{
		"message": "Oil quantity updated successfully",
	})
}
