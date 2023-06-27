package models

type DeliveryOrder struct {
	ID                int64  `gorm:"primary_key;auto_increment" json:"id"`
	Recipient         string `json:"recipient"`
	TravelOrderID     int64  `gorm:"not null" json:"travel_order_id"`
	CustomerLocation  string `json:"customer_location"`
	WarehouseLocation string `json:"warehouse_location"`
	OilId             int64  `gorm:"not null" json:"oil_id"`
	DeliveredQuantity int64  `json:"delivered_quantity"`
	WarehouseQuantity int64  `json:"warehouse_quantity"`
}
