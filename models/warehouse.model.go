package models

type Warehouse struct {
	Id         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"not null" json:"name"`
	Location   string `gorm:"not null" json:"location"`
	ProvinceId int    `gorm:"not null" json:"province_id"`
	CityId     int    `gorm:"not null" json:"city_id"`
	Storage    []Storage
}
