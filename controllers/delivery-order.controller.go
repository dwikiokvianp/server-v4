package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetDeliveryOrderByTravelOrderId(c *gin.Context) {

	var travelOrder models.TravelOrder

	if err := config.DB.Where("id = ?", c.Param("id")).
		First(&travelOrder).Error; err != nil {
		c.JSON(404, gin.H{"error": "Record not found!"})
		return
	}

	var deliveryOrder models.DeliveryOrder

	if err := config.DB.Where("travel_order_id = ?", c.Param("id")).
		First(&deliveryOrder).Error; err != nil {
		c.JSON(404, gin.H{"error": "Record not found!"})
		return
	}

	var deliveryOrderRecipient []models.DeliveryOrderRecipientDetail
	if err := config.DB.Where("delivery_order_id = ?", deliveryOrder.ID).
		Joins("Transaction.User.Company").
		Joins("Transaction.Officer").
		Joins("Transaction.Province").
		Joins("Transaction.City").
		Find(&deliveryOrderRecipient).Error; err != nil {
		c.JSON(404, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(200, gin.H{
		"delivery_order": deliveryOrder,
		"recipient":      deliveryOrderRecipient,
		"travel_order":   travelOrder,
	})

}

func GetDeliveryOrders(c *gin.Context) {
	var deliveryOrders []models.DeliveryOrderRecipientDetail
	config.DB.
		Joins("Transaction.User").
		Joins("Transaction.User.Company").
		Find(&deliveryOrders)

	c.JSON(200, deliveryOrders)
}
