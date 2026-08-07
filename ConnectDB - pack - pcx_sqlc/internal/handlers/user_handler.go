package handlers

import (
	"hoc-gin/internal/db/sqlc"
	"hoc-gin/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserResponse struct {
	UserID    int32  `json:"user_id"`
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (uh *UserHandler) GetUserByUuid(ctx *gin.Context) {
	uuidParam := ctx.Param("uuid")
	userUuid, err := uuid.Parse(uuidParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid uuid",
		})
		return
	}
	user, err := uh.repo.Find(ctx, userUuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	response := UserResponse{
		UserID:    user.UserID,
		Uuid:      user.Uuid.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02"),
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user": response,
	})
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input sqlc.CreateUserParams
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}
	user, err := uh.repo.Create(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	response := UserResponse{
		UserID:    user.UserID,
		Uuid:      user.Uuid.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02"),
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"user": response,
	})
}
