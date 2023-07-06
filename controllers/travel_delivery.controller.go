package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"server-v2/utils"
	"strconv"
)

func CreateTravelOrder(c *gin.Context) {
	travelDeliveryInput := models.TravelDeliveryInput{}

	if err := c.ShouldBindJSON(&travelDeliveryInput); err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to bind gagal",
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
			"message": "Failed to created travel order",
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

	go func() {
		for _, recipientDetail := range travelDeliveryInput.RecipientDetail {
			deliveryOrderRecipientDetail := models.DeliveryOrderRecipientDetail{
				DeliveryOrderID: deliveryOrder.ID,
				UserId:          recipientDetail.UserId,
				Quantity:        recipientDetail.Quantity,
				ProvinceId:      recipientDetail.ProvinceId,
				CityId:          recipientDetail.CityId,
				OilId:           travelDeliveryInput.OilId,
				TransactionID:   recipientDetail.TransactionID,
			}

			user := models.User{}

			fmt.Println("recipientDetail.UserId", recipientDetail.UserId)

			if err := config.DB.Where("id = ?", recipientDetail.UserId).First(&user).Error; err != nil {
				c.JSON(400, gin.H{
					"message": err.Error(),
				})
				return
			}

			if err := config.DB.Create(&deliveryOrderRecipientDetail).Error; err != nil {
				c.JSON(400, gin.H{
					"message": err.Error(),
				})
				return
			}

			qrData := recipientDetail.TransactionID
			qrFile, err := utils.GenerateQRCode(strconv.FormatInt(qrData, 10))
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			key := fmt.Sprintf("qrcodes/%v.png", qrData)
			qrURL, err := utils.UploadToS3(qrFile, key)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			email := user.Email
			err = utils.SendEmail("dwikiokvianp1999@gmail.com", "Qr Code Transaction", "Body", qrFile)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			fmt.Println("berhasil kirim email ke ", email)

			myTransaction := models.Transaction{}
			fmt.Println(recipientDetail.TransactionID)

			if err := config.DB.Where("id = ?", recipientDetail.TransactionID).First(&myTransaction).Error; err != nil {
				c.JSON(400, gin.H{
					"message": err.Error(),
				})
				return
			}

			myTransaction.QrCodeUrl = qrURL

			if err := config.DB.Save(&myTransaction).Error; err != nil {
				c.JSON(400, gin.H{
					"message": err.Error(),
				})
				return
			}
		}
	}()

	go func() {
		for _, warehouseDetail := range travelDeliveryInput.WarehouseDetail {
			deliveryOrderWarehouseDetail := models.DeliveryOrderWarehouseDetail{
				DeliveryOrderID: deliveryOrder.ID,
				WarehouseID:     warehouseDetail.WarehouseID,
				Quantity:        warehouseDetail.Quantity,
				StorageID:       warehouseDetail.StorageID,
			}

			if err := config.DB.Create(&deliveryOrderWarehouseDetail).Error; err != nil {
				c.JSON(400, gin.H{
					"message": err.Error(),
				})
				return
			}
		}
	}()

	c.JSON(200, gin.H{
		"message": "Success create travel order",
	})
}

func GetTravelOrder(c *gin.Context) {
	var travelOrder []models.TravelOrder
	if err := config.DB.
		Find(&travelOrder).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get travel order",
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
	if err := config.DB.Where("id = ?", id).First(&travelOrder).Error; err != nil {
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
