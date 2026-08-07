package repository

import (
	"hoc-gin/internal/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	Find(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
}
