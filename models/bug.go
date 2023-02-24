package models

type Bug struct {
	ID    uint   `json:"id" gorm:"primaryKey" gorm:"unique;unique_index"`
	Email string `json:"email"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// TODO: Rate limit the fuck out of this endpoint (throw back)
