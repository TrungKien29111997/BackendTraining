package handlers

import (
	"hoc-gin/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (uh *UserHandler) GetUserByUuid(ctx *gin.Context) {
	uh.repo.Find()
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get user by uuid",
	})
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	uh.repo.Create()
	ctx.JSON(http.StatusOK, gin.H{
		"message": "create user",
	})
}
