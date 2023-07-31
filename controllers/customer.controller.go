package controllers

import (
	"fmt"
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
func GetCustomerById(c *gin.Context) {
	var customer models.Customer

	customerId := c.Param("id")
	fmt.Print(customerId)

	if err := config.DB.
		Where("user_id = ?", customerId).
		Preload("User").
		Preload("Company").
		Find(&customer).Error; err != nil {
		c.JSON(404, gin.H{"message": "Record not found!"})
		return
	}

	c.JSON(200, customer)
}
