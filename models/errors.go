package models

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Error struct {
	Message    string      `json:"message"`
	CreatedAt  time.Time   `json:"created_at"`
	User       string      `json:"user"`
	IP         string      `json:"ip"`
	Submission []byte      `json:"submission"`
	Claims     jwt.Claims  `json:"claims"`
	Session    interface{} `json:"session"`
}
