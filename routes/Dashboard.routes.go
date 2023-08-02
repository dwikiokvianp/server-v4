package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func DashboardRoutes(c *gin.Engine) {
	dashboardRoutes := c.Group("")
	{
		dashboardRoutes.GET("/SummaryFuelOutData", controllers.SummaryFuelOutData)
		dashboardRoutes.GET("/TransactionSummaryData", controllers.TransactionSummaryData)
		dashboardRoutes.GET("/SummaryJettyData", controllers.SummaryJettyData)
		dashboardRoutes.GET("/SummarySPOBData", controllers.SummarySPOBData)
		dashboardRoutes.GET("/SummaryDeliveryData", controllers.SummaryDeliveryData)
	}
}
