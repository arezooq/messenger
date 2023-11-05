package handlers

import (
	"net/http"

	"github.com/arezooq/hex-messanger/internal/core/domain"
	"github.com/arezooq/hex-messanger/internal/core/services"
	"github.com/gin-gonic/gin"
)

type HTTPHandlerUser struct {
	svc	services.UserService
}

func NewHTTPHandlerUser(UserService services.UserService) *HTTPHandlerUser {
	return &HTTPHandlerUser{
		svc: UserService,
	}
}

func (h *HTTPHandlerUser) RegisterUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	err := h.svc.RegisterUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "New user created successfully",
	})
}

func (h *HTTPHandlerUser) GetOneUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := h.svc.GetOneUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (h *HTTPHandlerUser) GetAllUsers(ctx *gin.Context) {
	users, err := h.svc.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, users)
}

func (h *HTTPHandlerUser) LoginUser(ctx *gin.Context) {
	var user domain.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	response, err := h.svc.LoginUser(user.Email, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":            	response.ID,
		"email":         	response.Email,
		"AccessToken":		response.AccessToken,
	})
}

func (h *HTTPHandlerUser) UpdateUser(ctx *gin.Context) {
	var user domain.User

	id := ctx.Param("id")
	
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	userUpdate, err := h.svc.UpdateUser(id, user.Email, user.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":			"User successful updated",
		"id":            	userUpdate.ID,
		"email":         	userUpdate.Email,
	})
}

func (h *HTTPHandlerUser) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.svc.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}