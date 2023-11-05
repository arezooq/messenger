package domain

import "time"

type Message struct {
	ID		string	`json:"id"`
	Body	string	`json:"body"`
	UserID	string	`json:"user_id"`
}

type User struct {
	ID			string		`json:"id"`
	Email		string		`json:"email" validate:"email, required"`
	Password	string		`json:"password" validate:"required, min=8"`
	Created_at	time.Time	`json:"created_at"`
	Updated_at	time.Time	`json:"updated_at"`
}