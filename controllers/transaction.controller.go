package controllers

import (
	"fmt"
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
		VehicleId:  inputTransaction.VehicleId,
		OfficerId:  inputTransaction.OfficerId,
		Status:     "pending",
		Date:       inputTransaction.Date,
		CityId:     inputTransaction.CityId,
		ProvinceId: inputTransaction.ProvinceId,
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
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
			OilID:         detail.OilID,
			Quantity:      detail.Quantity,
			TransactionID: int64(transaction.ID),
			StorageID:     detail.StorageId,
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

	key := fmt.Sprintf("qrcodes/%v", transaction.ID)
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
	)

	pageParam := c.Query("page")
	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	db := config.DB

	var count int64
	if err := db.Model(&models.Transaction{}).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	offset := (page - 1) * pageSize

	err := db.Offset(offset).Limit(pageSize).
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

func GetByIdTransaction(c *gin.Context) {
	var transaction models.Transaction
	id := c.Param("id")
	err := config.DB.
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
	id := c.Param("id")
	fmt.Println(id)
	err := config.DB.Preload("Vehicle.VehicleType").Preload("Oil").Preload("User.Detail").Where("user_id = ?", id).Find(&transaction).Error
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

	c.JSON(200, gin.H{
		"data": transaction,
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

	today := time.Now().Format("2006-01-02")

	name := c.Query("name")

	if name != "" {
		err := config.DB.
			Preload("Vehicle.VehicleType").
			Preload("User.Company").Preload("Oil").
			Preload("User.Detail").
			Joins("JOIN users ON users.id = transactions.user_id").
			Where("users.username = ?", name).
			Where("date = ?", today).Find(&transaction).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
	} else {
		err := config.DB.Preload("Vehicle.VehicleType").Preload("User.Company").Preload("Oil").
			Preload("User.Detail").
			Where("date = ?", today).Find(&transaction).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
	}

	if len(transaction) == 0 {
		c.JSON(404, gin.H{
			"message": "Transaction Today Not Found",
		})
	}

	c.JSON(200, gin.H{
		"data": transaction,
	})

}
