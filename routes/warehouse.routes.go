package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func WarehouseRoutes(c *gin.Engine) {
	warehouseGroup := c.Group("/warehouses")
	{
		warehouseGroup.GET("", controllers.GetWarehouse)
		warehouseGroup.GET("/:id", controllers.GetWarehouseById)
		warehouseGroup.GET("/storage/:id", controllers.GetStorageWarehouse)
	}
}
