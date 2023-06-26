package models

type Officer struct {
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}
