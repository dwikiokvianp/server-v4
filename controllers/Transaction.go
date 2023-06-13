package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"strconv"
)

func GetAllTransactions(c *gin.Context) {
	var transactions []models.Transaction

	config.DB.Preload("User").Preload("Vehicle").Preload("Oil").Find(&transactions)

	c.JSON(200, gin.H{
		"data": transactions,
	})
}

func GetByIdTransaction(c *gin.Context) {
	var transaction models.Transaction
	id := c.Param("id")
	err := config.DB.Preload("Vehicle.VehicleType").Preload("User.Detail").Find(&transaction, id).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"data": transaction,
	})
}

func CreateTransactions(c *gin.Context) {
	var inputTransaction models.TransactionInput

	fmt.Printf("ini merupakan sampe")
	if err := c.ShouldBindJSON(&inputTransaction); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	userId := c.Param("id")

	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	transaction := models.Transaction{
		UserId:    intUserId,
		VehicleId: inputTransaction.VehicleId,
		OilId:     inputTransaction.OilId,
		QrCodeUrl: inputTransaction.QrCodeUrl,
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success create transaction",
	})

}
