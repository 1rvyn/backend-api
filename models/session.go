package models

// redis session string of key = userID and value = token

type Session struct {
	Role      int    `json:"role"`
	AccountID int    `json:"id" gorm:"primaryKey"`
	Token     string `json:"token"`
	Agent     string `json:"agent"`
}
