package models

type Customer struct {
	Id               int `json:"id" gorm:"primary_key"`
	CustomerDetailId int `json:"customer_detail_id"`
	CustomerTypeId   int `json:"customer_type_id"`
}

type CustomerDetail struct {
	Id         int `json:"id" gorm:"primary_key"`
	CustomerId int `json:"customer_id"`
}

type CustomerType struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
