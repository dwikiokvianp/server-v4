package controllers

import (
	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"strconv"
)

func CreateHandover(c *gin.Context) {
	var handover models.Handover

	if err := c.ShouldBindJSON(&handover); err != nil {
		c.JSON(500, gin.H{
			"error": "Error bind in the create handover",
		})
		return
	}

	id := c.MustGet("id").(string)
	idInt, _ := strconv.Atoi(id)

	handover.OfficerId = idInt
	handover.Status = "pending"

	if err := config.DB.Create(&handover).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Error here in the create handover",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": "success create handover",
	})
}

func GetHandoverById(c *gin.Context) {
	id := c.MustGet("id").(string)
	idInt, _ := strconv.Atoi(id)

	var handover models.Handover
	var handoverResponse models.HandoverResponse

	if err := config.DB.
		Where("officer_id = ?", idInt).
		Preload("WorkerBefore").
		Preload("WorkerAfter").
		Joins("Officer").
		Order("id desc").
		First(&handover).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Error here in the get handover by id",
		})
		return
	}

	err := dto.Map(&handoverResponse, &handover)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error here in the get handover by id",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": handoverResponse,
	})
}
