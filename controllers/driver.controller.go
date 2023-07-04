package controllers

import (
	"fmt"
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
	var travelOrders []models.TravelOrder
	fmt.Println(c.Param("id"))
	if err := config.DB.Where("driver_id = ?", c.Param("id")).
		Preload("DeliveryOrderRecipientDetail").
		Preload("DeliveryOrderWarehouseDetail").
		Find(&travelOrders).Error; err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": travelOrders,
	})
}
