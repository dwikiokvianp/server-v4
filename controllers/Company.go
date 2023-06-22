package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetAllCompany(c *gin.Context) {
	var company []models.Company

	err := config.DB.Find(&company).Error

	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error",
			"error":   err,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
		"data":    company,
	})

}
