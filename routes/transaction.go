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
		transactionGroup.GET("", controllers.GetAllTransactions)
		transactionGroup.GET("/:id", controllers.GetByIdTransaction)
		transactionGroup.GET("/user/:id", controllers.GetTransactionByUserId)
		transactionGroup.GET("/summary", GetTodayTransactionsHandler)
		transactionGroup.GET("/today", GetTodayTransactions)
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
		createdAt := time.Unix(transaction.CreatedAt, 0)
		if createdAt.Before(now) {
			orderDone++
			orderToday++
		} else if createdAt.After(now) {
			if transaction.DeliveryTime > endOfToday.Unix() && transaction.DeliveryTime < startOfTomorrow.Unix() {
				orderTomorrow++
			}
		}

		if transaction.Quantity > 0 {
			oilIn += transaction.Quantity
		} else {
			oilOut += transaction.Quantity
		}
	}

	response := gin.H{
		"order_done":     orderDone,
		"order_today":    orderToday,
		"order_tomorrow": orderTomorrow,
		"oil_in":         oilIn,
		"oil_out":        oilOut,
	}

	c.JSON(http.StatusOK, response)
}


func GetTodayTransactions(c *gin.Context) {
	transactions, err := controllers.GetTodayTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Membuat slice untuk menyimpan objek transaksi dalam format yang diinginkan
	response := make([]gin.H, 0)
	today := time.Now().Format("02-01-2006")
	for _, transaction := range transactions {
		createdAt := time.Unix(transaction.CreatedAt, 0).Format("02-01-2006")
		if createdAt == today {
			response = append(response, gin.H{
				"id":       transaction.ID,
				"name":     transaction.User.Username,
				"phone":    transaction.User.Phone,
				"date":     createdAt,
				"quantity": transaction.Quantity,
				"status":   transaction.Status,
			})
		}
	}

	c.JSON(http.StatusOK, response)
}


