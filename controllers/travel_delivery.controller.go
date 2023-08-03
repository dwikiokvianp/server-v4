package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"server-v2/config"
	"server-v2/models"
	"strconv"
)

func CreateTravelOrder(c *gin.Context) {
	travelDeliveryInput := models.TravelDeliveryInput{}

	if err := c.ShouldBindJSON(&travelDeliveryInput); err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to bind json",
		})
		return
	}

	var totalQuantity int64
	for _, recipientDetail := range travelDeliveryInput.RecipientDetail {
		totalQuantity += recipientDetail.Quantity
	}

	if totalQuantity > travelDeliveryInput.Quantity {
		c.JSON(400, gin.H{
			"message": "Total quantity of recipient detail is greater than quantity of travel order",
		})
		return
	}

	fmt.Println(travelDeliveryInput.DepartureDate)
	travelOrder := models.TravelOrder{
		DriverId:       travelDeliveryInput.DriverId,
		PickupLocation: travelDeliveryInput.PickupLocation,
		DepartureDate:  travelDeliveryInput.DepartureDate,
		Message:        travelDeliveryInput.Message,
		OfficerId:      travelDeliveryInput.OfficerId,
		VehicleId:      travelDeliveryInput.VehicleId,
		Quantity:       travelDeliveryInput.Quantity,
		Status:         "received",
	}

	if err := config.DB.Create(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	deliveryOrder := models.DeliveryOrder{
		TravelOrderID: travelOrder.ID,
		OilId:         travelDeliveryInput.OilId,
	}

	if err := config.DB.Create(&deliveryOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, recipientDetail := range travelDeliveryInput.RecipientDetail {
		deliveryOrderRecipientDetail := models.DeliveryOrderRecipientDetail{
			TransactionDeliveryID: recipientDetail.TransactionDeliveryID,
			DeliveryOrderID:       deliveryOrder.ID,
			Quantity:              recipientDetail.Quantity,
			ProvinceId:            recipientDetail.ProvinceId,
			CityId:                recipientDetail.CityId,
			OilId:                 travelDeliveryInput.OilId,
		}

		if err := config.DB.Create(&deliveryOrderRecipientDetail).Error; err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "Success create travel order",
	})

}

func GetTravelOrder(c *gin.Context) {
	var travelOrder []models.TravelOrder
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(400, gin.H{
			"message": "Invalid page parameter",
		})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		c.JSON(400, gin.H{
			"message": "Invalid limit parameter",
		})
		return
	}

	offset := (pageInt - 1) * limitInt

	var totalRecords int64
	if err := config.DB.Model(&models.TravelOrder{}).Count(&totalRecords).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get total record count",
		})
		return
	}

	totalPage := int(math.Ceil(float64(totalRecords) / float64(limitInt)))

	if err := config.DB.
		Preload("Driver.User").
		Preload("Vehicle.VehicleIdentifier").
		Preload("Vehicle.VehicleType").
		Offset(offset).
		Limit(limitInt).
		Find(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data":        travelOrder,
		"totalCount":  totalRecords,
		"currentPage": pageInt,
		"pageSize":    limitInt,
		"totalPage":   totalPage,
	})
}

func GetTravelOrderByUser(c *gin.Context) {
	var travelOrder []models.TravelOrder

	userId := c.MustGet("id")

	fmt.Print(userId)

	if err := config.DB.Where("driver_id = ?", userId).
		Find(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get travel order",
		})
	}

	if travelOrder == nil {
		c.JSON(400, gin.H{
			"message": "Travel order not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": travelOrder,
	})

}

func GetTravelOrderById(c *gin.Context) {
	id := c.Param("id")
	var travelOrder models.TravelOrder
	if err := config.DB.Where("id = ?", id).
		Preload("Driver.User").
		Preload("Vehicle.VehicleIdentifier").
		Preload("Vehicle.VehicleType").
		Preload("DeliveryOrderRecipientDetail.DeliveryOrder").
		Preload("DeliveryOrderRecipientDetail.TransactionDelivery.Transaction.Customer.User").
		Preload("DeliveryOrderRecipientDetail.TransactionDelivery.Transaction.Customer.Company").
		First(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Travel order not found",
		})
		return
	}

	travelOrder.Status = "read"
	config.DB.Save(&travelOrder)

	c.JSON(200, gin.H{
		"data": travelOrder,
	})
}
func UpdateBatchStatus(c *gin.Context) {
	id := c.Param("id")
	var travelOrder []models.TravelOrder
	if err := config.DB.Where("id = ?", id).
		Preload("DeliveryOrderRecipientDetail.Transaction.User.Company").
		Preload("DeliveryOrderRecipientDetail.Transaction.TransactionDetail").
		Preload("DeliveryOrderRecipientDetail.Transaction.Status.Status").
		Preload("Driver").
		Find(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Travel order not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Success update status",
	})
}

func UpdateStatusTravel(c *gin.Context) {
	id := c.Param("id")
	var travelOrder models.TravelOrder
	if err := config.DB.Where("id = ?", id).First(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Travel order not found",
		})
		return
	}

	if travelOrder.Status == "read" {
		travelOrder.Status = "ongoing"
	} else if travelOrder.Status == "ongoing" {
		travelOrder.Status = "delivering"
	} else if travelOrder.Status == "delivering" {
		travelOrder.Status = "delivered"
	}
	config.DB.Save(&travelOrder)

	c.JSON(200, gin.H{
		"data": travelOrder,
	})
}
