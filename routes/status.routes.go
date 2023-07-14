package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func StatusRoutes(c *gin.Engine) {
	statusGroup := c.Group("/status")
	{
		statusGroup.GET("", controllers.GetAllStatus)
	}
}
