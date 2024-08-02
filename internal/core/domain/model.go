package domain

import "time"

type Message struct {
	Id        string    `json:"_id" bson:"_id"`
	Body      string    `json:"body" bson:"body"`
	UserId    string    `json:"user_id" bson:"user_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type User struct {
	Id        string    `json:"_id" bson:"_id"`
	Email     string    `json:"email" bson:"email" validate:"email, required"`
	Password  string    `json:"-" bson:"password" validate:"required, min=8"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
