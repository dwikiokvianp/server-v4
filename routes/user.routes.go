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
		userGroup.POST("", controllers.CreateUser)
		userGroup.PUT("/balance-credit", controllers.UpdateBalanceAndCredit)
	}
}
