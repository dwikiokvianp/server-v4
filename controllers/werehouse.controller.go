package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetWarehouse(c *gin.Context) {
	warehouse := models.Warehouse{}
	if err := config.DB.
		Preload("WarehouseDetail.Storage").
		Find(&warehouse).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, warehouse)
}

func GetWarehouseById(c *gin.Context) {
	id := c.Param("id")
	warehouse := models.Warehouse{}
	if err := config.DB.
		Preload("WarehouseDetail.Storage").
		Where("id = ?", id).
		Find(&warehouse).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if warehouse.Id == 0 {
		c.JSON(404, gin.H{"message": "warehouse not found"})
		return
	}

	c.JSON(200, warehouse)
}
