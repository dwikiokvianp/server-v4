package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func DeliveryOrderRoutes(c *gin.Engine) {
	deliveryGroupRoutes := c.Group("/delivery-order")
	{
		deliveryGroupRoutes.GET("", controllers.GetDeliveryOrders)
		deliveryGroupRoutes.GET("/:id", controllers.GetDeliveryOrderByTravelOrderId)
	}
}
