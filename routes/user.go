package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func UserRoutes(c *gin.Engine) {
	userGroup := c.Group("/users")
	{
		userGroup.GET("", controllers.GetAllUser)
		userGroup.GET("/:id", controllers.GetUserById)
	}
}
