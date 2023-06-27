package models

type Storage struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	QuantityTank int64  `json:"quantity_tank"`
}
