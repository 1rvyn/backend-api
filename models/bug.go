package models

type Bug struct {
	Email string `json:"email" gorm:"primaryKey" gorm:"unique;unique_index" `
	Title string `json:"title"`
	Body  string `json:"body"`
}
