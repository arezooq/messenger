package services

import (
	"github.com/google/uuid"
	"messenger/internal/adapters/repositories"
	"messenger/internal/core/domain"
	"messenger/internal/core/ports"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) RegisterUser(user domain.User) error {
	user.Id = uuid.New().String()
	return u.repo.RegisterUser(user)
}

func (u *UserService) GetOneUser(id string) (*domain.User, error) {
	return u.repo.GetOneUser(id)
}

func (u *UserService) GetAllUsers() ([]*domain.User, error) {
	return u.repo.GetAllUsers()
}

func (u *UserService) UpdateUser(id, email, password string) (*domain.User, error) {
	return u.repo.UpdateUser(id, email, password)
}
func (u *UserService) DeleteUser(id string) error {
	return u.repo.DeleteUser(id)
}

func (u *UserService) LoginUser(email, password string) (*repositories.LoginResponse, error) {
	return u.repo.LoginUser(email, password)
}
