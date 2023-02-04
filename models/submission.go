package models

import (
	"time"
)

// TODO: 1. add a field for the submission type
// TODO: 2. add a field for the submission status
// TODO: 3. add a field for who submitted the submission (account ID/name)

// TODO: 4. unqiue index for what we search by
type Submission struct {
	ID         uint `json:"id" gorm:"primaryKey"`
	CreatedBy  uint `json:"created_by"`
	CreatedAt  time.Time
	Code       string `json:"code" gorm:"unique"`
	Cookie     string `json:"cookie"`
	Email      string `json:"email"`
	IP         string `json:"ip"`
	Successout string `json:"successout"`
	Errorout   string `json:"errorout"`
	//MetaData string `json:"meta_data"`
}
