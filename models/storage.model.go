package models

type Storage struct {
	ID                int64  `json:"id"`
	WarehouseDetailID uint64 `gorm:"not null" json:"warehouse_detail_id"`
	Name              string `json:"name"`
	Quantity          int64  `json:"quantity"`
	OilID             uint64 `gorm:"not null" json:"oil_id"`
}
