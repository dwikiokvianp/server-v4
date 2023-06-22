package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func SummaryRoutes(c *gin.Engine) {
	summaryGroup := c.Group("/summary")
	{
		summaryGroup.GET("", controllers.GetSummary)
	}
}
