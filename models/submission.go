package models

import (
	"time"
)

type Submission struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedBy uint `json:"created_by"`
	CreatedAt time.Time
	Code      string `json:"code" gorm:"unique"`
	Language  string `json:"language"`
	Status    string `json:"status" gorm:"default:pending"`
	UserID    string `json:"user-id" gorm:"foreignKey:UserID"`
	IP        string `json:"ip"`

	//MetaData string `json:"meta_data"`
}
