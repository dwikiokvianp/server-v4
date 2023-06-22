package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetSummary(c *gin.Context) {
	type Summary struct {
		OrderDone    int64
		OrderPending int64
		OrderToday   int64
	}

	var summary Summary
	today := config.DB.Model(&models.Transaction{}).Where("status = ?", "done").Count(&summary.OrderDone)
	config.DB.Model(&models.Transaction{}).Where("status = ?", "pending").Count(&summary.OrderPending)

	summary.OrderToday = today.RowsAffected

	c.JSON(200, gin.H{
		"summary": summary,
	})

}
