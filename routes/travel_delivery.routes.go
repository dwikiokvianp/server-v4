package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func TravelDeliveryRoutes(c *gin.Engine) {
	travelDeliveryGroup := c.Group("travel-delivery")
	{
		travelDeliveryGroup.POST("", controllers.CreateTravelOrder)
		travelDeliveryGroup.GET("", controllers.GetTravelOrder)
		travelDeliveryGroup.GET("/:id", controllers.GetTravelOrderById)
		travelDeliveryGroup.GET("/user", controllers.GetTravelOrderByUser)
		travelDeliveryGroup.PATCH("/status/batch/:id", controllers.UpdateBatchStatus)
		travelDeliveryGroup.PATCH("/:id", controllers.UpdateStatusTravel)

	}
}
