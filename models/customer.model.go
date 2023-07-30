package models

type Customer struct {
	Id             int     `json:"id" gorm:"primary_key"`
	UserId         int     `json:"user_id"`
	User           User    `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerTypeId int     `json:"customer_type_id"`
	DetailId       int     `json:"detail_id"`
	Detail         Detail  `json:"detail" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Company        Company `json:"company" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CompanyID      int     `json:"company_id"`
	CreatedAt      int64   `gorm:"autoCreateTime" json:"created_at"`
	Phone          string  `json:"phone"`
}

type CustomerInput struct {
	Id             int    `json:"id" gorm:"primary_key"`
	UserId         int    `json:"user_id"`
	CustomerTypeId int    `json:"customer_type_id"`
	DetailId       int    `json:"detail_id"`
	CompanyID      int    `json:"company_id"`
	CreatedAt      int64  `gorm:"autoCreateTime" json:"created_at"`
	Phone          string `json:"phone"`
}

type CustomerType struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
