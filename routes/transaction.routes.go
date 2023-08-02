package routes

import (
	"net/http"
	"server-v2/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(router *gin.Engine) {
	transactionGroup := router.Group("/transactions")
	{
		transactionGroup.POST("/:id", controllers.CreateTransactions)
		transactionGroup.GET("/active", controllers.GetTransactionDeliveryActive)
		transactionGroup.GET("", controllers.GetAllTransactions)
		transactionGroup.GET("/delivery/:id", controllers.GetTransactionDelivery)
		transactionGroup.GET("/delivery/detail/:id", controllers.GetTransactionDeliveryById)
		transactionGroup.GET("/:id", controllers.GetByIdTransaction)
		transactionGroup.PATCH("", controllers.UpdateTransactionBatch)
		transactionGroup.PATCH("/:id", controllers.UpdateTransaction)
		transactionGroup.GET("/user/:id", controllers.GetTransactionByUserId)
		transactionGroup.GET("/summary", GetTodayTransactionsHandler)
		transactionGroup.GET("/tomorrow", GetTomorrowTransactions)
		transactionGroup.GET("/today", GetTodayTransactions)
		transactionGroup.GET("/v2/toda", controllers.GetTodayV2Transaction)
		transactionGroup.PATCH("/status/:id", controllers.UpdateStatusTransactions)
		transactionGroup.POST("/postpone/:id", controllers.PostponeTransaction)
		transactionGroup.PATCH("/type/:id", controllers.UpdateTransactionType)
		transactionGroup.GET("/mantap/halo/jiwa", controllers.TestEndpoint)
		transactionGroup.GET("/user/transaction", controllers.GetUserTransaction)
	}
}

func GetTodayTransactionsHandler(c *gin.Context) {
	transactions, err := controllers.GetTodayTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderDone := 0
	orderToday := 0
	orderTomorrow := 0
	oilIn := 0
	oilOut := 0

	now := time.Now()
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	startOfTomorrow := endOfToday.Add(24 * time.Hour)

	for _, transaction := range transactions {
		if transaction.Date.Before(now) {
			orderDone++
			orderToday++
		} else if transaction.Date.After(endOfToday) && transaction.Date.Before(startOfTomorrow) {
			orderTomorrow++
		}
	}

	response := gin.H{
		"order_done":  orderDone,
		"order_today": orderToday,
		"oil_in":      oilIn,
		"oil_out":     oilOut,
	}

	c.JSON(http.StatusOK, response)
}

const defaultLimit = 10

func GetTomorrowTransactions(c *gin.Context) {

}

func GetTodayTransactions(c *gin.Context) {
	//offsetStr := c.Query("offset")
	//limitStr := c.Query("limit")
	//
	//offset, err := strconv.Atoi(offsetStr)
	//if err != nil || offset < 0 {
	//	offset = 0
	//}
	//
	//limit, err := strconv.Atoi(limitStr)
	//if err != nil || limit < 1 {
	//	limit = defaultLimit
	//}
	//
	//transactions, err := controllers.GetTodayTransactions()
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//username := c.Query("username")
	//
	//response := make([]gin.H, 0)
	//today := time.Now().Format("02-01-2006")
	//for _, transaction := range transactions[offset : offset+limit] {
	//	transactionDate := transaction.Date.Format("02-01-2006")
	//	if transactionDate == today && (username == "" || transaction.User.Username == username) {
	//		response = append(response, gin.H{
	//			"id":     transaction.ID,
	//			"name":   transaction.User.Username,
	//			"date":   transactionDate,
	//			"status": transaction.StatusId,
	//		})
	//	}
	//}
	//
	//c.JSON(http.StatusOK, response)
}
