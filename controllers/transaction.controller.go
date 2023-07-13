package controllers

import (
	"fmt"
	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
	"server-v2/config"
	"server-v2/models"
	"server-v2/utils"
	"strconv"
	"time"
)

func CreateTransactions(c *gin.Context) {

	var inputTransaction models.TransactionInput

	if err := c.ShouldBindJSON(&inputTransaction); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	userId := c.Param("id")

	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	transaction := models.Transaction{
		UserId:     intUserId,
		Email:      inputTransaction.Email,
		OfficerId:  inputTransaction.OfficerId,
		StatusId:   inputTransaction.StatusId,
		Date:       inputTransaction.Date,
		CityId:     inputTransaction.CityId,
		ProvinceId: inputTransaction.ProvinceId,
		DriverId:   inputTransaction.DriverId,
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Error here in the create transaction",
		})
		return
	}

	var transactionDetails []models.TransactionDetail

	if inputTransaction.TransactionDetail == nil {
		c.JSON(400, gin.H{
			"message": "transaction detail is required",
		})
		return
	}

	for _, detail := range inputTransaction.TransactionDetail {
		transactionDetail := models.TransactionDetail{
			Quantity:      detail.Quantity,
			TransactionID: int64(transaction.ID),
			StorageID:     detail.StorageId,
			OilID:         detail.OilID,
		}
		transactionDetails = append(transactionDetails, transactionDetail)
	}

	if err := config.DB.Create(&transactionDetails).Error; err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	qrData := strconv.Itoa(int(transaction.ID))
	qrFile, err := utils.GenerateQRCode(qrData)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	key := fmt.Sprintf("qrcodes/%v.png", transaction.ID)
	qrURL, err := utils.UploadToS3(qrFile, key)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	email := inputTransaction.Email
	subject := "QR Code Transaction"
	body := `
	<html>
	<head>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f2f2f2;
				padding: 20px;
			}
	
			h1 {
				color: #333333;
				font-size: 24px;
				font-weight: bold;
				margin-bottom: 20px;
			}
	
			p {
				color: #666666;
				font-size: 16px;
				line-height: 1.5;
				margin-bottom: 10px;
			}
	
			.qr-code {
				display: block;
				text-align: center;
				margin-bottom: 20px;
			}
	
			.qr-code img {
				max-width: 200px;
				height: auto;
			}
		</style>
	</head>
	<body>
		<p>Tunjukkan QR kode ini kepada petugas untuk mendapatkan layanan.</p>
	</body>
	</html>
	`

	go func() {
		err = utils.SendEmail(email, subject, body, qrFile)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}()

	transaction.QrCodeUrl = qrURL
	if err := config.DB.Save(&transaction).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success create transaction",
	})
}

func GetAllTransactions(c *gin.Context) {
	var (
		transactions []models.Transaction
		pageSize     = 10
		page         = 1
		status       = "pending"
	)

	pageParam := c.Query("page")
	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	limitParam := c.Query("limit")
	if limitParam != "" {
		pageSize, _ = strconv.Atoi(limitParam)
	}

	statusParam := c.Query("status")
	if statusParam != "" {
		status = statusParam
	}

	db := config.DB

	var count int64
	if err := db.Model(&models.Transaction{}).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	offset := (page - 1) * pageSize

	err := db.Offset(offset).Limit(pageSize).
		Where("status = ?", status).
		Joins("User.Company").
		Joins("User").
		Joins("Vehicle.VehicleType").
		Joins("Officer").
		Joins("Province").
		Joins("City").
		Preload("TransactionDetail").
		Order("updated_at desc").
		Find(&transactions).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	totalPages := (int(count) + pageSize - 1) / pageSize

	c.JSON(200, gin.H{
		"data":     transactions,
		"page":     page,
		"pageSize": pageSize,
		"total":    totalPages,
	})
}

func UpdateTransactionBatch(c *gin.Context) {
	type IdToUpdate struct {
		Id        uint64 `json:"id"`
		VehicleId *int   `json:"vehicle_id"`
		DriverId  int    `json:"driver_id"`
	}

	var ids []IdToUpdate
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, id := range ids {
		transaction := models.Transaction{}
		if err := config.DB.Find(&transaction, id.Id).Error; err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		status := c.Query("status")

		statusInt, err := strconv.Atoi(status)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Please provide status",
			})
			return
		}

		transaction.VehicleId = id.VehicleId
		transaction.DriverId = id.DriverId
		transaction.StatusId = statusInt
		if err := config.DB.Save(&transaction).Error; err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "success update transaction",
	})

}

func GetByIdTransaction(c *gin.Context) {
	var transaction models.Transaction
	id := c.Param("id")
	err := config.DB.
		Preload("TransactionDetail.Oil").
		Joins("Vehicle.VehicleType").
		Joins("User.Role").
		Joins("Officer").
		Joins("User.Detail").
		Joins("User.Company").
		Find(&transaction, id).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": transaction,
	})
}

func GetTransactionByUserId(c *gin.Context) {
	var transaction []models.Transaction
	var transactionResponse []models.TransactionResponse
	id := c.Param("id")
	fmt.Println(id)
	err := config.DB.
		Preload("TransactionDetail.Oil").
		Where("user_id = ?", id).
		Find(&transaction).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	if len(transaction) == 0 {
		c.JSON(404, gin.H{
			"message": "Not Found",
		})
		return
	}

	err = dto.Map(&transactionResponse, &transaction)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": transactionResponse,
	})
}

func GetTodayTransactions() ([]models.Transaction, error) {
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	var transactions []models.Transaction
	err := config.DB.
		Preload("User").
		Where("date BETWEEN ? AND ?", startOfToday, endOfToday).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func GetTomorrowTransactions() ([]models.Transaction, error) {
	now := time.Now()
	startOfTomorrow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
	endOfTomorrow := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Add(24 * time.Hour)

	var transactions []models.Transaction
	err := config.DB.
		Preload("User").
		Where("date BETWEEN ? AND ?", startOfTomorrow.Format("2006-01-02T15:04:05Z07:00"), endOfTomorrow.Format("2006-01-02T15:04:05Z07:00")).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func GetTodayV2Transaction(c *gin.Context) {
	var transaction []models.Transaction

	today := time.Now().UTC().Format("2006-01-02")

	name := c.Query("name")

	if name != "" {
		err := config.DB.
			Preload("Vehicle.VehicleType").
			Preload("User.Company").
			Preload("User.Detail").
			Joins("JOIN users ON users.id = transactions.user_id").
			Where("users.username = ?", name).
			Where("date = ?", today).
			Find(&transaction).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		err := config.DB.Preload("Vehicle.VehicleType").
			Preload("User.Company").
			Preload("User.Detail").
			Where("date = ?", today).
			Find(&transaction).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	if len(transaction) == 0 {
		c.JSON(404, gin.H{
			"message": "Transaction Today Not Found",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": transaction,
	})
}

func UpdateStatusTransactions(c *gin.Context) {
	transaction := models.Transaction{}

	id := c.Param("id")
	fmt.Println(id)
	err := config.DB.Find(&transaction, id).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if transaction.ID == 0 {
		c.JSON(404, gin.H{
			"message": "Transaction Not Found",
		})
		return
	}

	status := c.Query("status")
	if status == "" {
		c.JSON(400, gin.H{
			"message": "Status is required",
		})
		return
	}

	statusQuery := c.Query("status")
	if statusQuery == "" {
		c.JSON(400, gin.H{
			"message": "Status is required",
		})
		return
	}

	statusInt, err := strconv.Atoi(statusQuery)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Status must be integer",
		})
		return
	}

	transaction.StatusId = statusInt
	if err := config.DB.Save(&transaction).Error; err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Success update status to %s transaction with id %s", status, id),
	})
}
