package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"messenger/internal/core/domain"
	"messenger/internal/core/services"
)

type HTTPHandlerUser struct {
	svc services.UserService
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

func (h *HTTPHandlerUser) GetAllUsersByExportData(ctx *gin.Context) {
	users, err := h.svc.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	fmt.Println("users ::", users)

	file := excelize.NewFile()
	fmt.Println("file ::", file)

	sheetName := "Sheet1"
	sheetIndex, err := file.NewSheet(sheetName)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create a new sheet",
		})
		return
	}

	headers := []string{"ID", "Email", "Password", "Created_at", "Updated_at"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := file.SetCellValue(sheetName, cell, header); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}
	}

	for i, item := range users {
		row := i + 2
		values := []interface{}{item.Id, item.Email, item.Password, item.CreatedAt, item.UpdatedAt}
		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			if err := file.SetCellValue(sheetName, cell, value); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"Error": err.Error(),
				})
				return
			}
		}
	}

	file.SetActiveSheet(sheetIndex)
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=data.xlsx")
	if err := file.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}
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
		"id":          response.ID,
		"email":       response.Email,
		"AccessToken": response.AccessToken,
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
		"message": "User successful updated",
		"id":      userUpdate.Id,
		"email":   userUpdate.Email,
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
