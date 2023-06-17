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
		transactionGroup.GET("/tomorrow", GetTomorrowTransactions)
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
		if transaction.Date.Before(now) {
			orderDone++
			orderToday++
		} else if transaction.Date.After(endOfToday) && transaction.Date.Before(startOfTomorrow) {
			orderTomorrow++
		}

		if transaction.Quantity > 0 {
			oilIn += transaction.Quantity
		} else {
			oilOut += transaction.Quantity
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

func GetTomorrowTransactions(c *gin.Context) {
	transactions, err := controllers.GetTomorrowTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get query parameter "username" from the URL
	username := c.Query("username")

	// Create a slice to store transaction objects in the desired format
	response := make([]gin.H, 0)
	tomorrow := time.Now().Add(24 * time.Hour).Format("02-01-2006")
	for _, transaction := range transactions {
		transactionDate := transaction.Date.Format("02-01-2006")
		if transactionDate == tomorrow && (username == "" || transaction.User.Username == username) {
			response = append(response, gin.H{
				"id":       transaction.ID,
				"name":     transaction.User.Username,
				"phone":    transaction.User.Phone,
				"date":     transactionDate,
				"quantity": transaction.Quantity,
				"status":   transaction.Status,
			})
		}
	}

	c.JSON(http.StatusOK, response)
}


func GetTodayTransactions(c *gin.Context) {
	transactions, err := controllers.GetTodayTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get query parameter "username" from the URL
	username := c.Query("username")

	// Create a slice to store transaction objects in the desired format
	response := make([]gin.H, 0)
	today := time.Now().Format("02-01-2006")
	for _, transaction := range transactions {
		transactionDate := transaction.Date.Format("02-01-2006")
		if transactionDate == today && (username == "" || transaction.User.Username == username) {
			response = append(response, gin.H{
				"id":       transaction.ID,
				"name":     transaction.User.Username,
				"phone":    transaction.User.Phone,
				"date":     transactionDate,
				"quantity": transaction.Quantity,
				"status":   transaction.Status,
			})
		}
	}

	c.JSON(http.StatusOK, response)
}
