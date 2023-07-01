package models

import "time"

type Transaction struct {
	ID         uint64    `gorm:"primary_key:auto_increment" json:"id"`
	UserId     int       `gorm:"not null" json:"user_id"`
	Email      string    `gorm:"not null" json:"email"`
	User       User      `gorm:"foreignkey:UserId"`
	QrCodeUrl  string    `gorm:"not null" json:"qr_code_url"`
	VehicleId  int       `gorm:"not null" json:"vehicle_id"`
	Vehicle    Vehicle   `gorm:"foreignkey:VehicleId"`
	CreatedAt  int64     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  int64     `gorm:"autoUpdateTime" json:"updated_at"`
	OfficerId  int       `gorm:"not null" json:"officer_id"`
	Officer    Officer   `gorm:"foreignkey:OfficerId"`
	Status     string    `gorm:"not null" json:"status"`
	ProvinceId int       `gorm:"not null" json:"province_id"`
	Province   Province  `gorm:"foreignkey:ProvinceId"`
	CityId     int       `gorm:"not null" json:"city_id"`
	City       City      `gorm:"foreignkey:CityId"`
	Date       time.Time `json:"date"`
	DriverId   int       `gorm:"not null" json:"driver_id"`
}

type TransactionInput struct {
	VehicleId         int                      `gorm:"not null" json:"vehicle_id"`
	Email             string                   `gorm:"not null" json:"email"`
	OfficerId         int                      `gorm:"not null" json:"officer_id"`
	QrCodeUrl         string                   `json:"qr_code_url"`
	Date              time.Time                `json:"date"`
	CityId            int                      `gorm:"not null" json:"city_id"`
	ProvinceId        int                      `gorm:"not null" json:"province_id"`
	TransactionDetail []TransactionDetailInput `json:"transaction_detail"`
	StorageId         int                      `json:"storage_id"`
}

type TransactionDetailInput struct {
	OilID         int64 `json:"oil_id" binding:"required"`
	Quantity      int64 `json:"quantity" binding:"required"`
	TransactionId int64 `json:"transaction_id" binding:"required"`
	StorageId     int64 `json:"storage_id" binding:"required"`
	DriverId      int64 `json:"driver_id" binding:"required"`
}

type TransactionDetailBatchInput struct {
	Detail []TransactionDetailInput `json:"transaction_detail"`
}
