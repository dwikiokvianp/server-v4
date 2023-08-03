package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func UpdateTransactionDelivery(c *gin.Context) {
	deliveryID := c.Param("id")

	var updateRequest struct {
		Status string `form:"delivery_status"`
	}
	if err := c.ShouldBindQuery(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deliveryIDInt64, err := strconv.ParseInt(deliveryID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery ID"})
		return
	}

	var delivery models.TransactionDelivery
	if err := config.DB.First(&delivery, deliveryIDInt64).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction delivery not found"})
		return
	}

	delivery.DeliveryStatus = updateRequest.Status

	if err := config.DB.Save(&delivery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction delivery"})
		return
	}

	go func() {
		transactionId := delivery.TransactionID
		fmt.Println(transactionId, "ini transaction id")
		var transaction models.Transaction

		if err := config.DB.First(&transaction, transactionId).Error; err != nil {
			return
		}

		var transactionDeliveries []models.TransactionDelivery
		if err := config.DB.Where("transaction_id = ?", transactionId).Find(&transactionDeliveries).Error; err != nil {
			return
		}

		counterDelivery := 0
		for _, transactionDelivery := range transactionDeliveries {
			if transactionDelivery.DeliveryStatus != "pending" {
				counterDelivery++
				return
			} else {
				return
			}
		}

		if counterDelivery == len(transactionDeliveries) {
			transaction.IsFinished = true
			config.DB.Save(&transaction)
		}

	}()

	c.JSON(200, gin.H{
		"message": "success update status",
	})
}
