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
		OfficerID:      0,
		PickupLocation: travelDeliveryInput.PickupLocation,
		DepartureDate:  travelDeliveryInput.DepartureDate,
		Message:        travelDeliveryInput.Message,
		Status:         travelDeliveryInput.Status,
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
			"error": "Failed to create delivery order",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Travel delivery and controller is already created",
	})
}
