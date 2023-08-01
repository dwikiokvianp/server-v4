package models

type DeliveryOrder struct {
	ID            int64       `gorm:"primary_key;auto_increment" json:"id"`
	TravelOrderID int64       `gorm:"not null" json:"travel_order_id"`
	TravelOrder   TravelOrder `gorm:"foreignkey:TravelOrderID"`
	OilId         int64       `gorm:"not null" json:"oil_id"`
}

type DeliveryOrderWarehouseDetail struct {
	ID              int64 `gorm:"primary_key;auto_increment" json:"id"`
	DeliveryOrderID int64 `gorm:"not null" json:"delivery_order_id"`
	WarehouseID     int64 `gorm:"not null" json:"warehouse_id"`
	StorageID       int64 `gorm:"not null" json:"storage_id"`
	Quantity        int64 `json:"quantity"`
}

type DeliveryOrderRecipientDetail struct {
	ID                int64         `gorm:"primary_key;auto_increment" json:"id"`
	DeliveryOrderID   int64         `gorm:"not null" json:"delivery_order_id"`
	DeliveryOrder     DeliveryOrder `gorm:"foreignkey:DeliveryOrderID"`
	TransactionID     int64         `gorm:"not null" json:"transaction_id"`
	DeliveredQuantity int           `json:"delivered_quantity"`
	Transaction       Transaction   `json:"transaction"`
	OilId             int64         `gorm:"not null" json:"oil_id"`
	Quantity          int64         `json:"quantity"`
	ProvinceId        int64         `gorm:"not null" json:"province_id"`
	CityId            int64         `gorm:"not null" json:"city_id"`
}
