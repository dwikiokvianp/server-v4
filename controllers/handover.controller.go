package controllers

import (
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
	var id = c.Param("id")
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
	id := c.Param("id")

	var handover models.Handover
	var workerBefore models.Officer
	var workerAfter models.Officer

	if err := config.DB.
		Where("officer_id = ?", id).
		Joins("Officer").
		First(&handover).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Error here in the get handover by id",
		})
		return
	}

	if err := config.DB.
		Where("id = ?", handover.WorkerBefore).
		First(&workerBefore).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Error here in the get handover by id",
		})
		return
	}

	if err := config.DB.
		Where("id = ?", handover.WorkerAfter).
		First(&workerAfter).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Error here in the get handover by id",
		})
		return
	}

	c.JSON(200, gin.H{
		"data":          handover,
		"worker_before": workerBefore,
		"worker_after":  workerAfter,
	})
}
