package models

import "time"

type ResponseAccount struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique;unique_index"`
	CreatedAt time.Time `json:"created_at"`
	UserRole  int       `gorm:"default:1"`
}
