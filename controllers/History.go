package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"time"
)

func GetTodayHistory(c *gin.Context) {
	var history []models.HistoryOut
	dateNow := time.Now()
	dateStart := time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, time.UTC)
	dateEnd := dateStart.Add(24 * time.Hour)

	config.DB.Where("date >= ? AND date < ?", dateStart, dateEnd).Preload("User.Role").Preload("User.Company").Preload("Oil").Find(&history)

	var count int64
	config.DB.Model(&models.HistoryOut{}).Where("date >= ? AND date < ?", dateStart, dateEnd).Count(&count)

	var totalQuantity float64
	config.DB.Model(&models.HistoryOut{}).Where("date >= ? AND date < ?", dateStart, dateEnd).Select("SUM(quantity)").Scan(&totalQuantity)

	c.JSON(200, gin.H{
		"data":          history,
		"total":         count,
		"totalQuantity": totalQuantity,
	})
}
