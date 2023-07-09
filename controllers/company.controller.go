package controllers

import (
	"fmt"
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

func CreateCompany(c *gin.Context) {
	company := models.Company{}

	err := c.ShouldBindJSON(&company)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error",
			"error":   err,
		})
		return
	}

	err = config.DB.Create(&company).Error
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error",
			"error":   err,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Success create company with name %s", company.CompanyName),
	})
}

func DeleteCompany(c *gin.Context) {
	var company models.Company

	err := config.DB.Where("id = ?", c.Param("id")).First(&company).Error
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error",
			"error":   err,
		})
		return
	}

	if company.CompanyName == "" {
		c.JSON(400, gin.H{
			"message": "Error",
			"error":   "Company not found",
		})
		return
	}

	err = config.DB.Delete(&company).Error
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error",
		})
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Success delete company with name %s", company.CompanyName),
	})

}

func UpdateCompany(c *gin.Context) {
	var company models.Company

	err := config.DB.Where("id = ?", c.Param("id")).First(&company).Error
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error",
			"error":   err,
		})
		return
	}

	if company.CompanyName == "" {
		c.JSON(400, gin.H{
			"message": "Error",
			"error":   "Company not found",
		})
		return
	}

	err = c.ShouldBindJSON(&company)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error",
			"error":   err,
		})
		return
	}
	err = config.DB.Save(&company).Error
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error",
		})
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Success update company with name %s", company.CompanyName),
	})
}
