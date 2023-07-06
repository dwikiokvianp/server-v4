package models

type Handover struct {
	Id           uint64     `json:"id"`
	StartShiftId uint64     `gorm:"foreignkey:StartShiftId" json:"start_shift_id"`
	StartShift   StartShift `gorm:"foreignKey:StartShiftId" json:"start_shift"`
	EndShiftId   uint64     `gorm:"foreignkey:EndShiftId" json:"end_shift_id"`
	EndShift     EndShift   `gorm:"foreignKey:EndShiftId" json:"end_shift"`
}

type StartShift struct {
	Id           uint64 `json:"id"`
	WorkerBefore uint64 `gorm:"foreignKey:OfficerId" json:"worker_before"`
	OfficerId    uint64 `json:"officer_id"`
	Condition    string `json:"condition"`
	LevelStorage string `json:"level_storage"`
}

type EndShift struct {
	Id           uint64 `json:"id"`
	WorkerAfter  uint64 `gorm:"foreignKey:OfficerId" json:"worker_after"`
	OfficerId    uint64 `json:"officer_id"`
	Condition    string `json:"condition"`
	LevelStorage string `json:"level_storage"`
}
