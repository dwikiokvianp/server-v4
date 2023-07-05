package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetDeliveryOrderById(c *gin.Context) {
	var deliveryOrder models.DeliveryOrder
	id := c.Params.ByName("id")

	err := config.DB.Where("id = ?", id).First(&deliveryOrder).Error
	if err != nil {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, deliveryOrder)
	}
}

func GetDeliveryOrders(c *gin.Context) {
	var deliveryOrders []models.DeliveryOrder
	config.DB.Find(&deliveryOrders)

	c.JSON(200, deliveryOrders)
}
