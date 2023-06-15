package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func OfficerRoutes(c *gin.Engine) {
	officerGroup := c.Group("/officer")
	{
		officerGroup.GET("/", controllers.GetOfficer)
		officerGroup.GET("/:id", controllers.GetOfficerById)
	}
}
