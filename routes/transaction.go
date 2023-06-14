package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func TransactionRoutes(c *gin.Engine) {
	transactionGroup := c.Group("/transactions")
	{
		transactionGroup.GET("", controllers.GetAllTransactions)
		transactionGroup.GET("/user/:id", controllers.GetTransactionByUserId)
		transactionGroup.POST("/:id", controllers.CreateTransactions)
		transactionGroup.GET("/:id", controllers.GetByIdTransaction)
	}
}
