package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Province struct {
	Id   uint `json:"id" gorm:"primaryKey;autoIncrement" json:"province_id"`
	Name string
}

type City struct {
	Id         int `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string
	ProvinceID uint `json:"province_id"`
}

type ProvinceData struct {
	Provinsi string   `json:"provinsi"`
	Kota     []string `json:"kota"`
}

func main() {
	failedLoadEnv := godotenv.Load("./.env.local")
	if failedLoadEnv != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set in the environment variables")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&Province{}, &City{})
	if err != nil {
		log.Fatalf("Failed to perform database migration: %v", err)
	}

	jsonData, err := ioutil.ReadFile("./seed/regions.json")
	if err != nil {
		log.Fatalf("Failed to read JSON data: %v", err)
	}

	var provinceData []ProvinceData
	err = json.Unmarshal(jsonData, &provinceData)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Seed data
	for _, data := range provinceData {
		province := Province{Name: data.Provinsi}
		err = db.Create(&province).Error
		if err != nil {
			log.Fatalf("Failed to create province: %v", err)
		}

		for _, city := range data.Kota {
			cityData := City{Name: city, ProvinceID: province.Id}
			err = db.Create(&cityData).Error
			if err != nil {
				log.Fatalf("Failed to create city: %v", err)
			}
		}
	}

	fmt.Println("Seeding completed successfully")
}
