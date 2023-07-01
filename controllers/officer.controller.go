package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetOfficer(c *gin.Context) {
	var officer []models.Officer

	err := config.DB.Select("id", "username", "email").Find(&officer).Error
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Cannot load Officer",
		})
	}

	if len(officer) <= 0 {
		c.JSON(404, gin.H{
			"message": "Officer Not Found",
		})
	}

	c.JSON(200, gin.H{
		"data": officer,
	})

}

func GetOfficerById(c *gin.Context) {
	var officer models.Officer
	id := c.Param("id")

	err := config.DB.Select("id", "username").Where("id = ?", id).First(&officer).Error
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Internal Server Error",
		})
	}

	if officer == (models.Officer{}) {
		c.JSON(404, gin.H{
			"message": "Cannot find Officer",
		})
	}

	c.JSON(200, gin.H{
		"data": officer,
	})
}

func CreateOfficer(c *gin.Context) {
	officer := models.Officer{}

	err := c.ShouldBindJSON(&officer)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid Request",
		})
	}

	err = config.DB.Create(&officer).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
	}

	c.JSON(200, gin.H{
		"message": "Officer successfully registered",
	})
}
