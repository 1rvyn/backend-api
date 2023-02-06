package models

import (
	"time"
)

type Submission struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	Code      []byte `json:"code" gorm:"unique"`
	Language  string `json:"language"`
	Status    string `json:"status" gorm:"default:pending"`
	UserID    string `json:"user-id"`
	IP        string `json:"ip"`
}
