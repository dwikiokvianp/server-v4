package models

type Handover struct {
	Id           int     `json:"id"`
	WorkerBefore int     `json:"worker_before"`
	WorkerAfter  int     `json:"worker_after"`
	OfficerId    int     `json:"officer_id"`
	Officer      Officer `json:"officer" gorm:"foreignKey:OfficerId"`
	Condition    string  `json:"condition"`
}
