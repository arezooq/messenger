package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arezooq/hex-messanger/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserMongoRepository struct {
	client *mongo.Client
	db string
	collection *mongo.Collection
}

func NewUserMongoRepository() *UserMongoRepository {

	err := godotenv.Load(".env")
	
	if err != nil {
		log.Fatal("Error loading file .env")
	}

	Mongodb := os.Getenv("MONGO_USER_URL")

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

	collection := client.Database("hex-user").Collection("users")

	return &UserMongoRepository{
		client: client,
		db: Mongodb,
		collection: collection,
	}

}

func (u *UserMongoRepository) RegisterUser(user domain.User) error {

	errUserExist := u.UserMongoExist(user.Email)

	if errUserExist != nil {
		return errors.New("user already exists")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.Password = string(hashedPassword)

	_, err := u.collection.InsertOne(context.Background(), user)
	if err != nil {
		return errors.New(fmt.Sprintf("user not saved: %v", err.Error()))
	}
	return nil
}

func (u *UserMongoRepository) GetOneUser(id string) (*domain.User, error) {
	user := &domain.User{}
	err := u.collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil,  errors.New(fmt.Sprintf("user not found: %v", err.Error()))
	}
	return user, nil
}

func (u *UserMongoRepository) GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User
	req, err := u.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil,  errors.New(fmt.Sprintf("users not found: %v", err.Error()))
	}

	defer req.Close(context.Background())
	for req.Next(context.Background()) {
		var user *domain.User
		if err := req.Decode(&user); err != nil {
			return nil,  errors.New(fmt.Sprintf("users not found: %v", err.Error()))
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserMongoRepository) LoginUser(email, password string) (*LoginResponse, error) {
	user := &domain.User{}
	
	err := u.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil,  errors.New(fmt.Sprintf("user not found: %v", err.Error()))
	}

	err = u.VerifyMongoPassword(user.Password, password)
	if err != nil {
		return nil, errors.New("user not exists")
	}

	err = godotenv.Load(".env")
	
	if err != nil {
		return nil, errors.New("Error loading file .env")
	}

	JWTSecret := os.Getenv("SECRET_JWT")
	

	accessToken, err := u.GenerateMongoAccessToken(user.ID, JWTSecret)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		ID:           user.ID,
		Email:        user.Email,
		AccessToken:  accessToken,
	}, nil
}

func (u *UserMongoRepository) UpdateUser(id, email, password string) (*domain.User, error) {
	var user domain.User

	err := u.collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil,  errors.New(fmt.Sprintf("user not found: %v", err.Error()))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("password not hashed: %v", err))
	}

	user.Password = string(hashedPassword)
	user.Email = email

	update := bson.M{"email": user.Email, "password": user.Password}
	result, err := u.collection.UpdateOne(context.Background(), bson.M{"id": id}, bson.M{"$set": update})

	if err != nil {
		return nil, errors.New("unable to update user :(")
	}

	var updatedUser domain.User
	if result.MatchedCount == 1 {
		err := u.collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&updatedUser)
		if err != nil {
			return nil, errors.New("unable to found updated user :(")
		}
	}

	return &updatedUser, nil

}

func (u *UserMongoRepository) DeleteUser(id string) error {

	result, err := u.collection.DeleteOne(context.Background(), bson.M{"id": id})

	if err != nil {
		return errors.New("unable to delete user :(")
	}

	if result.DeletedCount < 1 {
		return errors.New("User with specified ID not found!")
	}
	return nil
}


func (u *UserMongoRepository) VerifyMongoPassword(hashedPassword, password string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("password not matched")
	}
	return nil
}

func (u *UserMongoRepository) GenerateMongoAccessToken(userID, jwtSecret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour).UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (u *UserMongoRepository) UserMongoExist(email string) error{
	user := &domain.User{}
	err := u.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return errors.New("user already exists")
	}
	return nil
}
