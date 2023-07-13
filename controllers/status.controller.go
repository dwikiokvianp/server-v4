package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetAllStatus(c *gin.Context) {
	status := models.Status{}

	if err := config.DB.Find(&status).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, status)
}
