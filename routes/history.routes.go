package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func HistoryRoutes(c *gin.Engine) {
	var history = c.Group("/history")
	{
		history.GET("/today", controllers.GetTodayHistory)
		history.GET("/today/out", controllers.GetHistoryOutToday)
	}
}
