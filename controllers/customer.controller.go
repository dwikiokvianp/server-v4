package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetAllCustomer(c *gin.Context) {
	var customer []models.Customer

	customerType := c.Query("customer_type")

	if err := config.DB.
		Where("customer_type_id = ?", customerType).
		Preload("User").
		Preload("Company").
		Find(&customer).Error; err != nil {
		c.JSON(404, gin.H{"message": "Record not found!"})
		return
	}

	if len(customer) == 0 {
		c.JSON(404, gin.H{"message": "Customer not found!"})
		return
	}

	c.JSON(200, customer)
}
