package repositories

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/arezooq/hex-messanger/internal/core/domain"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type MessangerPostgresRepository struct {
	db	*gorm.DB
}

func NewMessangerPostgresRepository() *MessangerPostgresRepository {
	err := godotenv.Load(".env")
	
	if err != nil {
		log.Fatal("Error loading file .env")
	}

	connStr := os.Getenv("POSTGRES_MESSANGER_URL")

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&domain.Message{})

	return &MessangerPostgresRepository{
		db: db,
	}
}

func (m *MessangerPostgresRepository) CreateMessage(message domain.Message) error {
	req := m.db.Create(&message)
	if req.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("messages not saved: %v", req.Error))
	}
	return nil
}

func (m *MessangerPostgresRepository) GetOneMessage(id string) (*domain.Message, error) {
	message := &domain.Message{}
	req := m.db.First(&message, "id = ? ", id)
	if req.RowsAffected == 0 {
		return nil, errors.New(fmt.Sprintf("message not found: %v", req.Error))
	}
	return message, nil
}

func (m *MessangerPostgresRepository) GetAllMessages() ([]*domain.Message, error) {
	var messages []*domain.Message
	req := m.db.Find(&messages)
	if req.RowsAffected == 0 {
		return nil, errors.New(fmt.Sprintf("messages not found: %v", req.Error))
	}
	return messages, nil
}

func (m *MessangerPostgresRepository) UpdateMessage(id, body, user_id string) (*domain.Message, error) {
	var message domain.Message

	req := m.db.First(&message, "id = ? ", id)
	if req.RowsAffected == 0 {
		return nil, errors.New("message not found")
	}
	message.Body = body

	req = m.db.Model(&message).Where("id = ? AND user_id = ?", id, user_id).Update(message)
	if req.RowsAffected == 0 {
		return nil, errors.New("unable to update message :(")
	}

	return &message, nil

}

func (m *MessangerPostgresRepository) DeleteMessage(id, user_id string) error {
	message := &domain.Message{}
	req := m.db.Where("id = ? AND user_id = ?", id, user_id).Delete(&message)
	if req.RowsAffected == 0 {
		return errors.New("message not found")
	}

	return nil
}