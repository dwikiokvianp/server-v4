package models

type DeliveryOrder struct {
	ID                int64  `gorm:"primary_key;auto_increment" json:"id"`
	Recipient         string `gorm:"not null" json:"recipient"`
	Message           string `gorm:"not null" json:"message"`
	TravelOrderID     int64  `gorm:"not null" json:"travel_order_id"`
	CustomerLocation  string `gorm:"not null" json:"customer_location"`
	WarehouseLocation string `gorm:"not null" json:"warehouse_location"`
	OilId             int64  `gorm:"not null" json:"oil_id"`
	DeliveredQuantity int64  `gorm:"not null" json:"delivered_quantity"`
	WarehouseQuantity int64  `gorm:"not null" json:"warehouse_quantity"`
}
