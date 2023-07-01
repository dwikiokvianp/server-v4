package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func CreateTravelOrder(c *gin.Context) {
	travelDeliveryInput := models.TravelDeliveryInput{}

	if err := c.ShouldBindJSON(&travelDeliveryInput); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	travelOrder := models.TravelOrder{
		DriverId:       travelDeliveryInput.DriverId,
		PickupLocation: travelDeliveryInput.PickupLocation,
		DepartureDate:  travelDeliveryInput.DepartureDate,
		Message:        travelDeliveryInput.Message,
		Status:         "received",
	}

	if err := config.DB.Create(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to created travel order",
		})
		return
	}

	deliveryOrder := models.DeliveryOrder{
		Recipient:         travelDeliveryInput.Recipient,
		TravelOrderID:     travelOrder.ID,
		CustomerLocation:  travelDeliveryInput.CustomerLocation,
		WarehouseLocation: travelDeliveryInput.WarehouseLocation,
		OilId:             travelDeliveryInput.OilId,
		DeliveredQuantity: travelDeliveryInput.DeliveredQuantity,
		WarehouseQuantity: travelDeliveryInput.WarehouseQuantity,
	}

	if err := config.DB.Create(&deliveryOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Travel delivery and controller is successfully created",
	})
}

func GetTravelOrderById(c *gin.Context) {
	id := c.Param("id")
	var travelOrder models.TravelOrder
	if err := config.DB.Where("id = ?", id).First(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Travel order not found",
		})
		return
	}

	travelOrder.Status = "read"
	config.DB.Save(&travelOrder)

	c.JSON(200, gin.H{
		"data": travelOrder,
	})
}

func GetTravelOrders(c *gin.Context) {
	var travelOrders []models.TravelOrder
	id := c.Param("driver_id")
	if err := config.DB.Where("driver_id = ?", id).Find(&travelOrders).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to get travel orders",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": travelOrders,
	})
}

func UpdateStatusTravel(c *gin.Context) {
	id := c.Param("id")
	var travelOrder models.TravelOrder
	if err := config.DB.Where("id = ?", id).First(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Travel order not found",
		})
		return
	}

	if travelOrder.Status == "read" {
		travelOrder.Status = "ongoing"
	} else if travelOrder.Status == "ongoing" {
		travelOrder.Status = "delivering"
	} else if travelOrder.Status == "delivering" {
		travelOrder.Status = "delivered"
	}
	config.DB.Save(&travelOrder)

	c.JSON(200, gin.H{
		"data": travelOrder,
	})
}
