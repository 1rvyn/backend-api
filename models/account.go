package models

import "time"

type Account struct {
	ID                int       `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name"`
	Email             string    `json:"email" gorm:"unique;unique_index"`
	EncryptedPassword []byte    `json:"password"`
	CreatedAt         time.Time `json:"created_at"`
	UserRole          int       `gorm:"default:1"`
	EmailCode         int       `json:"email_code"`
	Verified          bool      `json:"email_verified"`
}
