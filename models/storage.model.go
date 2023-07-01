package models

type Storage struct {
	ID                int64  `json:"id"`
	WarehouseDetailID uint64 `gorm:"not null" json:"warehouse_detail_id"`
	Name              string `json:"name"`
	QuantityTank      int64  `json:"quantity_tank"`
}
