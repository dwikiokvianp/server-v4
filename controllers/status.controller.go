package controllers

import (
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"strconv"
)

func GetAllStatus(c *gin.Context) {
	var statusTypeMapping []models.StatusTypeMapping
	var statusTypeMappingResponse []models.StatusTypeMappingResponse

	statusTypeID := 1
	statusTypeIDQuery := c.Query("status_type_id")
	if statusTypeIDQuery != "" {
		intStatus, err := strconv.Atoi(statusTypeIDQuery)
		if err != nil {
			c.JSON(404, gin.H{"error": "status_type_id must be a number"})
			return
		}
		statusTypeID = intStatus
	}

	err := config.DB.Where("status_type_id = ?", statusTypeID).
		Joins("Status").
		Joins("StatusType").
		Find(&statusTypeMapping).Error
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	for _, item := range statusTypeMapping {
		statusTypeMappingResponse = append(statusTypeMappingResponse, models.StatusTypeMappingResponse{
			ID:         item.ID,
			Status:     item.Status.Name,
			StatusType: item.StatusType.Name,
		})
	}

	c.JSON(200, statusTypeMappingResponse)
}

func GetAllTypes(c *gin.Context) {
	var statusType []models.StatusType

	err := config.DB.Find(&statusType).Error
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, statusType)
}
