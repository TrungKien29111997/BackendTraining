package handlers

import (
	"hoc-gin/internal/models"
	"hoc-gin/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid uuid",
			})
		}
		return
	}
	var user *models.User
	if err := uh.repo.Find(id, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get user by uuid",
	})
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var user *models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}
	if err := uh.repo.Create(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": user,
	})
}
