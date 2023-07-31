package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetAllEmployee(c *gin.Context) {
	var employee []models.Employee

	role := c.Query("role_id")

	if err := config.DB.
		Where("role_id = ?", role).
		Preload("User").
		Preload("Role").
		Find(&employee).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}
	c.JSON(200, employee)
}
