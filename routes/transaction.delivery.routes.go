package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func UpdateTransactionRoutes(c *gin.Engine) {
	transactionDeliveryGroup := c.Group("/transaction-delivery")
	{
		transactionDeliveryGroup.PATCH("/:id", controllers.UpdateTransactionDelivery)
	}
}
