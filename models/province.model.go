package models

type Province struct {
	Id   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	City []City `json:"city" gorm:"foreignKey:ProvinceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
