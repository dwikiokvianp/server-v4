package models

type Company struct {
	Id             int    `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyName    string `json:"username"`
	Address        string `json:"password"`
	CompanyDetail  string `json:"company_detail"`
	CompanyZipCode int    `json:"company_zip_code"`
	CreatedAt      int64  `gorm:"autoCreateTime" json:"created_at"`
}
