package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetWarehouse(c *gin.Context) {
	warehouse := models.Warehouse{}
	if err := config.DB.Find(&warehouse).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, warehouse)
}
