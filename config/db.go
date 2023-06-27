package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"server-v2/models"
)

var DB *gorm.DB

func InitDatabase(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	migrateError := db.AutoMigrate(
		&models.Role{},
		&models.DeliveryOrder{},
		&models.TravelOrder{},
		&models.User{},
		&models.Oil{},
		&models.VehicleType{},
		&models.Vehicle{},
		&models.Transaction{},
		&models.TransactionDetail{},
		&models.Officer{},
		&models.Company{},
		&models.HistoryOut{},
		&models.HistoryIn{},
		&models.Detail{},
		&models.Storage{},
		&models.Proof{})
	if migrateError != nil {
		log.Fatalln(err)
		return nil
	}

	DB = db

	return nil

}
