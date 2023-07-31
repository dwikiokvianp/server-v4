package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func EmployeeRoutes(c *gin.Engine) {
	employeeRoutes := c.Group("/employees")
	{
		employeeRoutes.GET("", controllers.GetAllEmployee)
	}
}
