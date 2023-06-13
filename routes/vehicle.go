package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func VehicleRoutes(c *gin.Engine) {
	vehicleGroup := c.Group("/vehicle")
	{
		vehicleGroup.GET("/", controllers.GetAllVehicle)
		vehicleGroup.GET("/:id", controllers.GetVehicleById)
		vehicleGroup.POST("/", controllers.CreateVehicle)
	}
}
