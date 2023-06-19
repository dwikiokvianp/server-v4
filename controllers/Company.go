package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetAllCompany(c *gin.Context) {
	var company models.Company

	err := config.DB.Find(&company).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error when fetching all companies",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": company,
	})
}
