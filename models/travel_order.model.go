package models

import "time"

type TravelOrder struct {
	ID                           int64                          `gorm:"primary_key;auto_increment" json:"id"`
	DriverId                     int64                          `gorm:"not null" json:"driver_id"`
	PickupLocation               string                         `gorm:"not null" json:"pickup_location"`
	DepartureDate                time.Time                      `gorm:"not null" json:"departure_date"`
	Message                      string                         `gorm:"not null" json:"message"`
	Status                       string                         `gorm:"not null" json:"status"`
	OfficerId                    int64                          `gorm:"not null" json:"officer_id"`
	VehicleId                    int64                          `gorm:"not null" json:"vehicle_id"`
	Quantity                     int64                          `gorm:"not null" json:"quantity"`
	DeliveryOrderRecipientDetail []DeliveryOrderRecipientDetail `gorm:"foreignkey:delivery_order_id" json:"recipient_detail"`
	DeliveryOrderWarehouseDetail []DeliveryOrderWarehouseDetail `gorm:"foreignkey:delivery_order_id" json:"warehouse_detail"`
}

type TravelDeliveryInput struct {
	DriverId        int64                          `gorm:"not null" json:"driver_id" `
	OfficerId       int64                          `gorm:"not null" json:"officer_id" `
	Quantity        int64                          `gorm:"not null" json:"quantity" `
	PickupLocation  string                         `gorm:"not null" json:"pickup_location" `
	DepartureDate   time.Time                      `gorm:"not null" json:"departure_date" `
	Message         string                         `gorm:"not null" json:"message" `
	Status          string                         `gorm:"not null" json:"status" `
	OilId           int64                          `gorm:"not null" json:"oil_id" `
	VehicleId       int64                          `gorm:"not null" json:"vehicle_id"`
	RecipientDetail []DeliveryOrderRecipientDetail `json:"recipient_detail"`
	WarehouseDetail []DeliveryOrderWarehouseDetail `json:"warehouse_detail"`
}
