package models

import (
	"time"
)

type Error struct {
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	User       string    `json:"user"`
	IP         string    `json:"ip"`
	Submission []byte    `json:"submission"`
}
