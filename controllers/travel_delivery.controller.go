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
		OfficerId:      travelDeliveryInput.OfficerId,
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

	for _, recipientDetail := range travelDeliveryInput.RecipientDetail {
		deliveryOrderRecipientDetail := models.DeliveryOrderRecipientDetail{
			DeliveryOrderID: deliveryOrder.ID,
			UserId:          recipientDetail.UserId,
			Quantity:        recipientDetail.Quantity,
		}

		if err := config.DB.Create(&deliveryOrderRecipientDetail).Error; err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	for _, warehouseDetail := range travelDeliveryInput.WarehouseDetail {
		deliveryOrderWarehouseDetail := models.DeliveryOrderWarehouseDetail{
			DeliveryOrderID: deliveryOrder.ID,
			WarehouseID:     warehouseDetail.WarehouseID,
			Quantity:        warehouseDetail.Quantity,
		}

		if err := config.DB.Create(&deliveryOrderWarehouseDetail).Error; err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "Success create travel order",
	})
}

func GetTravelOrder(c *gin.Context) {
	var travelOrder []models.TravelOrder
	if err := config.DB.Find(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get travel order",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": travelOrder,
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

func UpdateWarehouseCustomerDetails(c *gin.Context) {
	var inputTravelOrder models.TravelOrder
	if err := c.ShouldBindJSON(&inputTravelOrder); err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to bind JSON",
		})
		return
	}

	var travelOrder models.TravelOrder
	if err := config.DB.Where("id = ?", inputTravelOrder.ID).First(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Travel order not found",
		})
		return
	}
}
