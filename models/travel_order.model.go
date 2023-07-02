package models

type TravelOrder struct {
	ID             int64  `gorm:"primary_key;auto_increment" json:"id"`
	DriverId       int64  `gorm:"not null" json:"driver_id"`
	PickupLocation string `gorm:"not null" json:"pickup_location"`
	DepartureDate  string `json:"departure"`
	Message        string `gorm:"not null" json:"message"`
	Status         string `gorm:"not null" json:"status"`
}

type TravelDeliveryInput struct {
	UserID            int64  `gorm:"not null" json:"user_id" binding:"required"`
	DriverId          int64  `gorm:"not null" json:"driver_id" binding:"required"`
	PickupLocation    string `gorm:"not null" json:"pickup_location" binding:"required"`
	DepartureDate     string `gorm:"not null" json:"departure_date" `
	Message           string `gorm:"not null" json:"message" binding:"required"`
	Status            string `gorm:"not null" json:"status" binding:"required"`
	CustomerLocation  string `json:"customer_location" binding:"required"`
	WarehouseLocation string `json:"warehouse_location" binding:"required"`
	OilId             int64  `gorm:"not null" json:"oil_id" binding:"required"`
	DeliveredQuantity int64  `json:"delivered_quantity" binding:"required"`
	WarehouseQuantity int64  `json:"warehouse_quantity" binding:"required"`
}
