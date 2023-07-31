package models

type Storage struct {
	ID          int64  `json:"id"`
	WarehouseID int    `gorm:"not null" json:"warehouse_id"`
	Name        string `json:"name"`
	Quantity    int64  `json:"quantity"`
	Capacity    int    `json:"capacity"`
	OilID       uint64 `gorm:"not null" json:"oil_id"`
}
