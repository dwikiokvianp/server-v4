package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func TravelDeliveryRoutes(c *gin.Engine) {
	travelDeliveryGroup := c.Group("travel-delivery")
	{
		travelDeliveryGroup.POST("", controllers.CreateTravelOrder)
	}
}
