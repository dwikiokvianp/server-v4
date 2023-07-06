package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetTodayHistory(c *gin.Context) {
	var historyOut []models.HistoryOut
	var historyIn []models.HistoryIn

	config.DB.Preload("User.Role").Preload("User.Company").Preload("Oil").Find(&historyOut)
	config.DB.Preload("Oil").Find(&historyIn)

	var countIn int64
	var countOut int64
	config.DB.Model(&models.HistoryOut{}).Count(&countOut)
	config.DB.Model(&models.HistoryIn{}).Count(&countIn)

	var totalQuantityOut float64
	var totalQuantityIn float64
	config.DB.Model(&models.HistoryOut{}).Select("SUM(quantity)").Scan(&totalQuantityOut)
	config.DB.Model(&models.HistoryIn{}).Select("SUM(quantity)").Scan(&totalQuantityIn)

	c.JSON(200, gin.H{
		"historyOut":       historyOut,
		"historyIn":        historyIn,
		"totalIn":          countIn,
		"totalOut":         countOut,
		"totalQuantityOut": totalQuantityOut,
		"totalQuantityIn":  totalQuantityIn,
	})
}

func GetHistoryOutToday(c *gin.Context) {
	var historyOut []models.HistoryOut

	if err := config.DB.
		Joins("Oil").
		Joins("User").
		Find(&historyOut).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": historyOut,
	})
}
