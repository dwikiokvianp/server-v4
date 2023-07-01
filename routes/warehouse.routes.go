package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func WarehouseRoutes(c *gin.Engine) {
	c.Group("/warehouse")
	{
		c.GET("", controllers.GetWarehouse)
	}
}
