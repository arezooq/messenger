package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arezooq/hex-messanger/internal/core/domain"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessangerMongoRepository struct {
	client *mongo.Client
	db string
	collection *mongo.Collection
}

func NewMessangerMongoRepository() *MessangerMongoRepository {
	err := godotenv.Load(".env")
	
	if err != nil {
		log.Fatal("Error loading file .env")
	}

	Mongodb := os.Getenv("MONGO_MESSANGER_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(Mongodb))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("hex-messanger").Collection("messages")

	return &MessangerMongoRepository{
		client: client,
		db: Mongodb,
		collection: collection,
	}

}

func (m *MessangerMongoRepository) CreateMessage(message domain.Message) error {
	_, err := m.collection.InsertOne(context.Background(), message)
	if err != nil {
		return errors.New(fmt.Sprintf("messages not saved: %v", err.Error()))
	}
	return nil
}

func (m *MessangerMongoRepository) GetOneMessage(id string) (*domain.Message, error) {
	message := &domain.Message{}
	err := m.collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&message)
	if err != nil {
		return nil,  errors.New(fmt.Sprintf("message not found: %v", err.Error()))
	}
	return message, nil
}

func (m *MessangerMongoRepository) GetAllMessages() ([]*domain.Message, error) {
	var messages []*domain.Message
	req, err := m.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil,  errors.New(fmt.Sprintf("messages not found: %v", err.Error()))
	}

	defer req.Close(context.Background())
	for req.Next(context.Background()) {
		var message *domain.Message
		if err := req.Decode(&message); err != nil {
			return nil,  errors.New(fmt.Sprintf("messages not found: %v", err.Error()))
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (m *MessangerMongoRepository) UpdateMessage(id, body, user_id string) (*domain.Message, error) {
	var message domain.Message

	filter := bson.M{"id": id, "userid": user_id}

	err := m.collection.FindOne(context.Background(), filter).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil,  errors.New(("message not found"))
		}
		return nil, err
	}

	message.Body = body

	update := bson.M{"$set":bson.M{"body": message.Body}}
	result, err := m.collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return nil, errors.New("unable to update message :(")
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("unable to found updated message :(")
	}

	return &message, nil

}

func (m *MessangerMongoRepository) DeleteMessage(id, user_id string) error {
	var message domain.Message

	filter := bson.M{"id": id, "userid": user_id}

	err := m.collection.FindOne(context.Background(), filter).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(("message not found"))
		}
		return err
	}

	result, err := m.collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return errors.New("unable to delete message :(")
	}

	if result.DeletedCount < 1 {
		return errors.New("Message with specified ID not found!")
	}
	return nil
}