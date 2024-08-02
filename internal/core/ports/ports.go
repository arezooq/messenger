package ports

import (
	"messenger/internal/adapters/repositories"
	"messenger/internal/core/domain"
)

type MessangerService interface {
	CreateMessage(userId string, message domain.Message) error
	GetOneMessage(id string) (*domain.Message, error)
	GetAllMessages() ([]*domain.Message, error)
	UpdateMessage(id, body, user_id string) (*domain.Message, error)
	DeleteMessage(id string) error
}

type UserService interface {
	RegisterUser(user domain.User) error
	GetOneUser(id string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	LoginUser(email, password string) (*repositories.LoginResponse, error)
	UpdateUser(id, email, password string) (*domain.User, error)
	DeleteUser(id string) error
}

type MessangerRepository interface {
	CreateMessage(message domain.Message) error
	GetOneMessage(id string) (*domain.Message, error)
	GetAllMessages() ([]*domain.Message, error)
	UpdateMessage(id, body, user_id string) (*domain.Message, error)
	DeleteMessage(id, user_id string) error
}

type UserRepository interface {
	RegisterUser(user domain.User) error
	GetOneUser(id string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	LoginUser(email, password string) (*repositories.LoginResponse, error)
	UpdateUser(id, email, password string) (*domain.User, error)
	DeleteUser(id string) error
}
