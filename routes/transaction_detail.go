package routes

import (
	"github.com/gin-gonic/gin"
	"server-v2/controllers"
)

func TransactionDetail(c *gin.Engine) {
	transactionDetailGroup := c.Group("/transaction_detail")
	{
		transactionDetailGroup.GET("", controllers.FindAllDetailTransaction)
		transactionDetailGroup.GET("/:id", controllers.FindDetailTransactionById)
	}
}
