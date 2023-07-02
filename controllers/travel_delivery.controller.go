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
			"message": "Failed to bind JSON",
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
			"message": "Failed to created travel order",
		})
		return
	}

	deliveryOrder := models.DeliveryOrder{
		TravelOrderID: travelOrder.ID,
		OilId:         travelDeliveryInput.OilId,
	}

	if err := config.DB.Create(&deliveryOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	deliveryOrderWarehouseDetail := models.DeliveryOrderWarehouseDetail{
		WarehouseID:     travelDeliveryInput.WarehouseID,
		Quantity:        travelDeliveryInput.WarehouseQuantity,
		DeliveryOrderID: deliveryOrder.ID,
	}

	if err := config.DB.Create(&deliveryOrderWarehouseDetail).Error; err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	deliveryOrderRecipientDetail := models.DeliveryOrderRecipientDetail{
		DeliveryOrderID: deliveryOrder.ID,
		UserId:          travelDeliveryInput.UserID,
		Quantity:        travelDeliveryInput.DeliveredQuantity,
	}

	if err := config.DB.Create(&deliveryOrderRecipientDetail).Error; err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Success create travel order",
	})
}

func GetTravelOrderById(c *gin.Context) {
	id := c.Param("id")
	var travelOrder models.TravelOrder
	if err := config.DB.Where("id = ?", id).First(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Travel order not found",
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
			"message": "Failed to get travel orders",
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
			"message": "Travel order not found",
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
