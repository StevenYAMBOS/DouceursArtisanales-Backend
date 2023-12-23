package models

import "time"

type Contact struct {
	ID       uint   `json:"id"`
	Username string `json:username`
	Email    string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Submitted_at time.Time `json:"submitted_at"`
}