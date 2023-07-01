package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func WarehouseRoutes(c *gin.Engine) {
	warehouseGroup := c.Group("/warehouse")
	{
		warehouseGroup.GET("", controllers.GetWarehouse)
	}
}
