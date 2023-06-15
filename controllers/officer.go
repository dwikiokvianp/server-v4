package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetOfficer(c *gin.Context) {
	var officer []models.Officer

	err := config.DB.Find(&officer).Error
	if err != nil {
		c.JSON(404, gin.H{
			"msg": "Not Found",
		})
	}

	c.JSON(200, gin.H{
		"data": officer,
	})

}

func GetOfficerById(c *gin.Context) {
	var officer models.Officer
	id := c.Param("id")

	err := config.DB.Where("id = ?", id).First(&officer).Error
	if err != nil {
		c.JSON(404, gin.H{
			"msg": "Not Found",
		})
	}

	c.JSON(200, gin.H{
		"data": officer,
	})
}
