package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
)

func GetAllVehicle(c *gin.Context) {
	var vehicle []models.Vehicle

	config.DB.Preload("VehicleType").Find(&vehicle)

	c.JSON(200, gin.H{
		"data": vehicle,
	})
}

func GetVehicleById(c *gin.Context) {
	var vehicle models.Vehicle
	id := c.Param("id")
	config.DB.Preload("VehicleType").Find(&vehicle, id)
	c.JSON(200, gin.H{
		"data": vehicle,
	})
}

func CreateVehicle(c *gin.Context) {
	var vehicle models.Vehicle

	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := config.DB.Create(&vehicle).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"message": "Vehicle created successfully!",
		"data":    vehicle,
	})

}
