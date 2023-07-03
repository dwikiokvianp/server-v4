package models

import "time"

type TravelOrder struct {
	ID             int64     `gorm:"primary_key;auto_increment" json:"id"`
	DriverId       int64     `gorm:"not null" json:"driver_id"`
	PickupLocation string    `gorm:"not null" json:"pickup_location"`
	DepartureDate  time.Time `gorm:"not null" json:"departure_date"`
	Message        string    `gorm:"not null" json:"message"`
	Status         string    `gorm:"not null" json:"status"`
	OfficerId      int64     `gorm:"not null" json:"officer_id"`
}

type TravelDeliveryInput struct {
	DriverId        int64                          `gorm:"not null" json:"driver_id" binding:"required"`
	OfficerId       int64                          `gorm:"not null" json:"officer_id" binding:"required"`
	PickupLocation  string                         `gorm:"not null" json:"pickup_location" binding:"required"`
	DepartureDate   time.Time                      `gorm:"not null" json:"departure_date" binding:"required"`
	Message         string                         `gorm:"not null" json:"message" binding:"required"`
	Status          string                         `gorm:"not null" json:"status" binding:"required"`
	OilId           int64                          `gorm:"not null" json:"oil_id" binding:"required"`
	RecipientDetail []DeliveryOrderRecipientDetail `json:"recipient_detail"`
	WarehouseDetail []DeliveryOrderWarehouseDetail `json:"warehouse_detail"`
}
