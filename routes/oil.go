package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func OilRoutes(c *gin.Engine) {
	oilGroup := c.Group("/oil")
	{
		oilGroup.GET("", controllers.GetAllOil)
	}
}
