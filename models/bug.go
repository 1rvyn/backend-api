package models

type Bug struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Title string `json:"title"`
	Body  string `json:"body" gorm:"primaryKey"`
}
