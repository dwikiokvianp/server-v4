package models

type DeliveryOrder struct {
	ID            int64 `gorm:"primary_key;auto_increment" json:"id"`
	TravelOrderID int64 `gorm:"not null" json:"travel_order_id"`
	OilId         int64 `gorm:"not null" json:"oil_id"`
}

type DeliveryOrderWarehouseDetail struct {
	ID              int64 `gorm:"primary_key;auto_increment" json:"id"`
	DeliveryOrderID int64 `gorm:"not null" json:"delivery_order_id"`
	WarehouseID     int64 `gorm:"not null" json:"warehouse_id"`
	StorageID       int64 `gorm:"not null" json:"storage_id"`
	Quantity        int64 `json:"quantity"`
}

type DeliveryOrderRecipientDetail struct {
	ID              int64  `gorm:"primary_key;auto_increment" json:"id"`
	DeliveryOrderID int64  `gorm:"not null" json:"delivery_order_id"`
	OilId           int64  `gorm:"not null" json:"oil_id"`
	UserId          int64  `gorm:"not null" json:"user_id"`
	Email           string `gorm:"not null" json:"email"`
	Quantity        int64  `json:"quantity"`
	ProvinceId      int64  `gorm:"not null" json:"province_id"`
	CityId          int64  `gorm:"not null" json:"city_id"`
}
