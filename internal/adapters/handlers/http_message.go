package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/arezooq/hex-messanger/internal/core/domain"
	"github.com/arezooq/hex-messanger/internal/core/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type HTTPHandlerMessanger struct {
	svcMessanger	services.MessangerService
}

func NewHTTPHandlerMessanger(MessangerService services.MessangerService) *HTTPHandlerMessanger {
	return &HTTPHandlerMessanger{
		svcMessanger: MessangerService,
	}
}

func (h *HTTPHandlerMessanger) CreateMessage(ctx *gin.Context) {
	var message domain.Message
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	err := godotenv.Load(".env")
	
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	JWTSecret := os.Getenv("SECRET_JWT")

	userID, err := ValidateToken(ctx.Request.Header.Get("Authorization"), JWTSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "user not authorization",
		})
		return
	}

	err = h.svcMessanger.CreateMessage(userID, message)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "New message created successfully",
	})
}

func (h *HTTPHandlerMessanger) GetOneMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	message, err := h.svcMessanger.GetOneMessage(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, message)
}

func (h *HTTPHandlerMessanger) GetAllMessages(ctx *gin.Context) {
	messages, err := h.svcMessanger.GetAllMessages()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, messages)
}

func (h *HTTPHandlerMessanger) UpdateMessage(ctx *gin.Context) {
	
	var message domain.Message

	id := ctx.Param("id")
	
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	err := godotenv.Load(".env")
	
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	JWTSecret := os.Getenv("SECRET_JWT")

	userID, err := ValidateToken(ctx.Request.Header.Get("Authorization"), JWTSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "user not authorization",
		})
		return
	}

	messageUpdate, err := h.svcMessanger.UpdateMessage(id, message.Body, userID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":			"Message successful updated",
		"id":            	messageUpdate.ID,
		"body":         	messageUpdate.Body,
		"user_id":			messageUpdate.UserID,
	})
}

func (h *HTTPHandlerMessanger) DeleteMessage(ctx *gin.Context) {
	id := ctx.Param("id")

	err := godotenv.Load(".env")
	
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	JWTSecret := os.Getenv("SECRET_JWT")

	userID, err := ValidateToken(ctx.Request.Header.Get("Authorization"), JWTSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "user not authorization",
		})
		return
	}

	err = h.svcMessanger.DeleteMessage(id, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Message deleted successfully",
	})
}

func ValidateToken(authHeader string, jwtSecret string) (string, error) {
	if authHeader == "" {
		return "", errors.New("token not found")
	}

	tokenString := authHeader[7:]

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("token not valid")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now().UTC()) {
		return "", errors.New("token has expired")
	}

	userID := claims.Subject

	return userID, nil
}