package models

type WarehouseDetail struct {
	ID          uint64 `gorm:"primary_key;auto_increment" json:"id"`
	WarehouseID uint64 `gorm:"not null" json:"warehouse_id"`
	Storage     []Storage
}
