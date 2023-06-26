package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func AuthRoutes(c *gin.Engine) {
	authGroup := c.Group("/auth")
	{
		authGroup.POST("/register", controllers.RegisterUser)
		authGroup.POST("/login", controllers.LoginUser)
	}
}
