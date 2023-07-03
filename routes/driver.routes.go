package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func DriverRoutes(c *gin.Engine) {
	driverGroup := c.Group("drivers")
	{
		driverGroup.GET("", controllers.GetDrivers)
		driverGroup.GET("/:id/transactions", controllers.GetTransactionByDriverId)
	}
}
