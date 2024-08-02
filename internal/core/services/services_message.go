package services

import (
	"github.com/google/uuid"
	"messenger/internal/core/domain"
	"messenger/internal/core/ports"
)

type MessangerService struct {
	repo ports.MessangerRepository
}

func NewMessangerService(repo ports.MessangerRepository) *MessangerService {
	return &MessangerService{
		repo: repo,
	}
}

func (m *MessangerService) CreateMessage(userId string, message domain.Message) error {
	message.Id = uuid.New().String()
	message.UserId = userId
	return m.repo.CreateMessage(message)
}

func (m *MessangerService) GetOneMessage(id string) (*domain.Message, error) {
	return m.repo.GetOneMessage(id)
}

func (m *MessangerService) GetAllMessages() ([]*domain.Message, error) {
	return m.repo.GetAllMessages()
}

func (m *MessangerService) UpdateMessage(id, body, user_id string) (*domain.Message, error) {
	return m.repo.UpdateMessage(id, body, user_id)
}

func (m *MessangerService) DeleteMessage(id, user_id string) error {
	return m.repo.DeleteMessage(id, user_id)
}
