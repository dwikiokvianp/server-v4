package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"server-v2/models"
	"os"
	"server-v2/config"
	"server-v2/routes"
	"time"
)

func main() {
	failedLoadEnv := godotenv.Load(".env")
	if failedLoadEnv != nil {
		log.Fatal("Error loading .env file")
	}

	databaseUrl, port := os.Getenv("DB_URL"), os.Getenv("PORT")
	config.InitDatabase(databaseUrl)

	server := gin.New()
	myCorsConfig := cors.DefaultConfig()
	myCorsConfig.AllowAllOrigins = true
	myCorsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	server.Use(cors.New(myCorsConfig))
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Check github action",
		})
	})
	routes.Routes(server)

		// Inisialisasi cron job
		c := cron.New()

		c.AddFunc("50 15 * * *", func() {
			// Query transaksi dengan status "WAITING CUSTOMER TO PICK UP" (ID: 3) pada hari ini
			var transactions []models.Transaction
			if err := config.DB.Where("status_id = ? AND DATE(date) = ?", 3, time.Now().Format("2006-01-02")).Find(&transactions).Error; err != nil {
				log.Println("Gagal mendapatkan transaksi yang akan ditunda:", err.Error())
				return
			}
	
			// Perbarui status transaksi menjadi "POSTPONED" (ID: 7)
			for _, transaction := range transactions {
				transaction.StatusId = 7
				if err := config.DB.Save(&transaction).Error; err != nil {
					log.Println("Gagal memperbarui status transaksi:", err.Error())
					continue
				}
	
				postponeHistory := models.PostponeHistory{
					TransactionID: int(transaction.ID),
					Reason:        "automatic postponed",
				}
				if err := config.DB.Create(&postponeHistory).Error; err != nil {
					log.Println("Gagal membuat entri PostponeHistory:", err.Error())
					continue
				}
			}

			log.Println("Cron job selesai: Transaksi berhasil diubah")
		})	

// Mulai cron job
c.Start()

	err := server.Run(port)
	if err != nil {
		log.Fatal("Error running server")
	}
}
