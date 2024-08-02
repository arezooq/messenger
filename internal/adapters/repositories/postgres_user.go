package repositories

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"messenger/internal/core/domain"
)

type LoginResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type UserPostgresRepository struct {
	db *gorm.DB
}

func NewUserPostgresRepository() *UserPostgresRepository {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading file .env")
	}

	connStr := os.Getenv("POSTGRES_USER_URL")

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&domain.User{})

	return &UserPostgresRepository{
		db: db,
	}
}

func (u *UserPostgresRepository) RegisterUser(user domain.User) error {

	errUserExist := u.UserExist(user.Email)
	if errUserExist != nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("password not hashed")
	}

	user.Password = string(hashedPassword)

	req := u.db.Create(&user)

	if req.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("user not saved: %v", req.Error))
	}
	return nil
}

func (u *UserPostgresRepository) GetOneUser(id string) (*domain.User, error) {
	user := &domain.User{}
	req := u.db.First(&user, "id = ? ", id)
	if req.RowsAffected == 0 {
		return nil, errors.New(fmt.Sprintf("user not found: %v", req.Error))
	}
	return user, nil
}

func (u *UserPostgresRepository) GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User
	req := u.db.Find(&users)
	if req.RowsAffected == 0 {
		return nil, errors.New(fmt.Sprintf("users not found: %v", req.Error))
	}
	return users, nil
}

func (u *UserPostgresRepository) LoginUser(email, password string) (*LoginResponse, error) {
	user := &domain.User{}

	req := u.db.First(&user, "email = ? ", email)
	if req.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	err := u.VerifyPassword(user.Password, password)
	if err != nil {
		return nil, errors.New("user not exists")
	}

	err = godotenv.Load(".env")

	if err != nil {
		return nil, errors.New("Error loading file .env")
	}

	JWTSecret := os.Getenv("SECRET_JWT")

	accessToken, err := u.GenerateAccessToken(user.Id, JWTSecret)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		ID:          user.Id,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil
}

func (u *UserPostgresRepository) UpdateUser(id, email, password string) (*domain.User, error) {
	var user domain.User

	req := u.db.First(&user, "id = ? ", id)
	if req.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("password not hashed: %v", err))
	}

	user.Password = string(hashedPassword)
	user.Email = email

	req = u.db.Model(&user).Where("id = ?", id).Update(user)
	if req.RowsAffected == 0 {
		return nil, errors.New("unable to update user :(")
	}

	return &user, nil

}

func (u *UserPostgresRepository) DeleteUser(id string) error {
	user := &domain.User{}
	req := u.db.Where("id = ?", id).Delete(&user)
	if req.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (u *UserPostgresRepository) VerifyPassword(hashedPassword, password string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("password not matched")
	}
	return nil
}

func (u *UserPostgresRepository) GenerateAccessToken(userID, jwtSecret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour).UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (u *UserPostgresRepository) UserExist(email string) error {
	user := &domain.User{}
	req := u.db.First(&user, "email = ? ", email)
	if req.RowsAffected != 0 {
		return errors.New(fmt.Sprintf("user already exists: %v", req.Error))
	}
	return nil
}
