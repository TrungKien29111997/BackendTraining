package handlers

import (
	"hoc-gin/internal/models"
	"hoc-gin/internal/repository"
	"net/http"
	"strconv"

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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid uuid",
		})
		return
	}
	uh.repo.Find(id)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get user by uuid",
	})
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}
	uh.repo.Create(user)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "create user",
	})
}
