package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"strconv"
)

func FindAllDetailTransaction(c *gin.Context) {
	var (
		transactions []models.TransactionDetail
		pageSize     = 10
		page         = 1
	)

	pageParam := c.Query("page")
	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	db := config.DB

	var count int64
	if err := db.Model(&models.TransactionDetail{}).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	offset := (page - 1) * pageSize

	if err := config.DB.Offset(offset).Limit(pageSize).Preload("Transaction").
		Preload("Transaction.Vehicle.VehicleType").
		Preload("Transaction.User.Role").
		Preload("Transaction.User.Company").
		Preload("Transaction.City").
		Preload("Transaction.Officer").
		Preload("Transaction.Province").
		Preload("Oil").Find(&transactions).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	totalPages := (int(count) + pageSize - 1) / pageSize
	c.JSON(200, gin.H{
		"data":     transactions,
		"page":     page,
		"pageSize": pageSize,
		"total":    totalPages,
	})
}

func FindDetailTransactionByTransactionId(c *gin.Context) {
	var transactionDetail []models.TransactionDetail

	if err := config.DB.Where("transaction_id = ?", c.Param("id")).
		Joins("Transaction.Vehicle.VehicleType").
		Joins("Transaction.User.Role").
		Joins("Transaction.User.Company").
		Joins("Transaction.Officer").
		Joins("Transaction.City").
		Joins("Transaction.Province").
		Joins("Oil").
		Find(&transactionDetail).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(transactionDetail) == 0 {
		c.JSON(404, gin.H{
			"message": fmt.Sprintf("Transaction with id %s not found", c.Param("id")),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": transactionDetail,
	})
}
