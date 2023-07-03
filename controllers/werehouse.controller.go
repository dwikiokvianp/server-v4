package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetWarehouse(c *gin.Context) {
	var warehouse []models.Warehouse
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

func GetStorageWarehouse(c *gin.Context) {
	var storageWarehouse []models.Storage

	if err := config.DB.
		Where("warehouse_detail_id = ?", c.Param("id")).
		Find(&storageWarehouse).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(storageWarehouse) == 0 {
		c.JSON(404, gin.H{"message": "storage not found"})
		return
	}

	c.JSON(200, gin.H{"data": storageWarehouse})
}
