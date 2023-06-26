package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func CityRoutes(c *gin.Engine) {
	cityGroup := c.Group("/city")
	{
		cityGroup.GET("", controllers.FindCity)
		cityGroup.GET("/:id", controllers.FindCityById)
	}
}
