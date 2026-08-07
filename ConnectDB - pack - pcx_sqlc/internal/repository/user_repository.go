package repository

import (
	"fmt"
	"hoc-gin/internal/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SQLUserRepository struct {
	db sqlc.Querier
}

func NewSQLUserRepository(DB sqlc.Querier) UserRepository {
	return &SQLUserRepository{
		db: DB,
	}
}

func (u *SQLUserRepository) Create(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := u.db.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("fail to create user: %w", err)
	}
	return user, nil
}

func (u *SQLUserRepository) Find(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := u.db.GetUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("fail to create user: %w", err)
	}
	return user, nil
}
