package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"time"
)

func GetTodayHistory(c *gin.Context) {
	var historyOut []models.HistoryOut
	var historyIn []models.HistoryIn
	dateNow := time.Now()
	dateStart := time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, time.UTC)
	dateEnd := dateStart.Add(24 * time.Hour)

	config.DB.Where("date >= ? AND date < ?", dateStart, dateEnd).Preload("User.Role").Preload("User.Company").Preload("Oil").Find(&historyOut)
	config.DB.Where("date >= ? AND date < ?", dateStart, dateEnd).Preload("Oil").Find(&historyIn)

	var countIn int64
	var countOut int64
	config.DB.Model(&models.HistoryOut{}).Where("date >= ? AND date < ?", dateStart, dateEnd).Count(&countOut)
	config.DB.Model(&models.HistoryIn{}).Where("date >= ? AND date < ?", dateStart, dateEnd).Count(&countIn)

	var totalQuantityOut float64
	var totalQuantityIn float64
	config.DB.Model(&models.HistoryOut{}).Where("date >= ? AND date < ?", dateStart, dateEnd).Select("SUM(quantity)").Scan(&totalQuantityOut)
	config.DB.Model(&models.HistoryIn{}).Where("date >= ? AND date < ?", dateStart, dateEnd).Select("SUM(quantity)").Scan(&totalQuantityIn)

	c.JSON(200, gin.H{
		"historyOut":       historyOut,
		"historyIn":        historyIn,
		"totalIn":          countIn,
		"totalOut":         countOut,
		"totalQuantityOut": totalQuantityOut,
		"totalQuantityIn":  totalQuantityIn,
	})
}
