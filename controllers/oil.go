package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetAllOil(c *gin.Context) {
	var oil []models.Oil

	config.DB.Find(&oil)

	c.JSON(200, gin.H{
		"data": oil,
	})
}
