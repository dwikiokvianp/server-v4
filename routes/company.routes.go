package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func CompanyRoutes(c *gin.Engine) {
	companyGroup := c.Group("/company")
	{
		companyGroup.GET("", controllers.GetAllCompany)
		companyGroup.POST("", controllers.CreateCompany)
		companyGroup.DELETE("/:id", controllers.DeleteCompany)
		companyGroup.PUT("/:id", controllers.UpdateCompany)
	}
}
