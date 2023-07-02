package models

type DeliveryOrder struct {
	ID                int64  `gorm:"primary_key;auto_increment" json:"id"`
	Recipient         string `json:"recipient"`
	TravelOrderID     int64  `gorm:"not null" json:"travel_order_id"`
	UserID            int64  `gorm:"not null" json:"user_id"`
	WarehouseID       int64  `gorm:"not null" json:"warehouse_id"`
	OilId             int64  `gorm:"not null" json:"oil_id"`
	DeliveredQuantity int64  `json:"delivered_quantity"`
	WarehouseQuantity int64  `json:"warehouse_quantity"`
}
