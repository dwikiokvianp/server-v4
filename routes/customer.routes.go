package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func CustomerRoutes(c *gin.Engine) {
	customerGroup := c.Group("/customers")
	{
		customerGroup.GET("", controllers.GetAllCustomer)
		customerGroup.GET("/:id", controllers.GetCustomerById)
	}
}
