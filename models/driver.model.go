package models

type Driver struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
