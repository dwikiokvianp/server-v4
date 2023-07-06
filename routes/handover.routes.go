package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func HandoverRoutes(c *gin.Engine) {
	handoverGroup := c.Group("handovers")
	{
		handoverGroup.POST("/:id", controllers.CreateHandover)
		handoverGroup.GET("/:id", controllers.GetHandoverById)
	}
}
