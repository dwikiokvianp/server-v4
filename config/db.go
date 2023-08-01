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

	migrateError := db.AutoMigrate(modelsToMigrate...)
	if migrateError != nil {
		log.Fatalln(err)
		return nil
	}

	DB = db
	return nil
}

func DropDatabase(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	dropError := db.Migrator().DropTable(modelsToMigrate...)
	if dropError != nil {
		log.Fatalln(err)
		return nil
	}

	return nil
}

var modelsToMigrate = []interface{}{
	&models.Role{},
	&models.DeliveryOrder{},
	&models.TravelOrder{},
	&models.Warehouse{},
	&models.User{},
	&models.Oil{},
	&models.Status{},
	&models.StatusType{},
	&models.StatusTypeMapping{},
	&models.VehicleType{},
	&models.Vehicle{},
	&models.Transaction{},
	&models.TransactionDetail{},
	&models.PostponeHistory{},
	&models.Officer{},
	&models.Company{},
	&models.HistoryOut{},
	&models.HistoryIn{},
	&models.Detail{},
	&models.Storage{},
	&models.Proof{},
	&models.DeliveryOrderRecipientDetail{},
	&models.DeliveryOrderWarehouseDetail{},
	&models.Handover{},
	&models.Customer{},
	&models.CustomerType{},
	&models.Employee{},
	&models.TransactionDelivery{},
}
