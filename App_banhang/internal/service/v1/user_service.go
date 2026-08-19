package v1service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo  repository.UserRepository
	cache *cache.RedisCacheService
}

func NewUserService(repo repository.UserRepository, redis *redis.Client) UserService {
	return &userService{
		repo:  repo,
		cache: cache.NewRedisCacheService(redis),
	}
}

func (us *userService) GetAllUsers(ctx *gin.Context, search, orderBy, sort string, page, limit int32) ([]sqlc.User, int32, error) {
	context := ctx.Request.Context()

	// Get Cache Redis
	cacheKey := us.generateCacheKey(search, orderBy, sort, page, limit)
	var cacheData struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}
	if err := us.cache.Get(cacheKey, &cacheData); err == nil && cacheData.Users != nil {
		return cacheData.Users, cacheData.Total, nil
	}

	if sort == "" {
		sort = "desc"
	}
	if orderBy == "" {
		orderBy = "user_created_at"
	}
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limitInt := utils.GetIntEnv("LIMIT_ITEM_PER_PAGE", 10)
		limit = int32(limitInt)
	}
	offset := (page - 1) * limit

	users, err := us.repo.GetAll(context, search, orderBy, sort, limit, offset)
	if err != nil {
		return nil, 0, utils.WrapError(err, "Failed to get users", utils.ErrCodeInternal)
	}
	total, err := us.repo.CountUsers(context, search)
	if err != nil {
		return nil, 0, utils.WrapError(err, "Failed to count users", utils.ErrCodeInternal)
	}

	// Create Cache Redis
	cacheData = struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}{
		Users: users,
		Total: int32(total),
	}
	us.cache.Set(cacheKey, cacheData, 5*time.Minute)
	return users, int32(total), nil
}

func (us *userService) CreateUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()
	input.UserEmail = utils.NormalizeString(input.UserEmail)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "Failed to hash password", utils.ErrCodeInternal)
	}
	input.UserPassword = string(hashedPassword)
	user, err := us.repo.Create(context, input)
	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "Failed to create user", utils.ErrCodeInternal)
	}
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}
	return user, nil
}

func (us *userService) GetUserByUUID(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	getUser, err := us.repo.GetByUUID(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.WrapError(err, "User not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "Failed to update user", utils.ErrCodeInternal)
	}
	return getUser, nil
}

func (us *userService) UpdateUser(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	if input.UserPassword != nil && *input.UserPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, utils.WrapError(err, "Failed to hash password", utils.ErrCodeInternal)
		}
		hashed := string(hashedPassword)
		input.UserPassword = &hashed
	}

	updatedUser, err := us.repo.Update(context, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.WrapError(err, "User not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "Failed to update user", utils.ErrCodeInternal)
	}
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}
	return updatedUser, nil
}

func (us *userService) SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	deletedUser, err := us.repo.SoftDelete(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.WrapError(err, "User not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "Failed to delete user", utils.ErrCodeInternal)
	}
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}
	return deletedUser, nil
}

func (us *userService) HardDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	deletedUser, err := us.repo.HardDelete(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.WrapError(err, "User not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "Failed to delete user", utils.ErrCodeInternal)
	}
	return deletedUser, nil
}

func (us *userService) RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	restoreUser, err := us.repo.Restore(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.WrapError(err, "User not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "Failed to restore user", utils.ErrCodeInternal)
	}
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}
	return restoreUser, nil
}

func (us *userService) generateCacheKey(search, orderBy, sort string, page, limit int32) string {
	search = strings.TrimSpace(search)
	if search == "" {
		search = "none"
	}
	orderBy = strings.TrimSpace(orderBy)
	if orderBy == "" {
		orderBy = "user_created_at"
	}
	sort = strings.ToLower(strings.TrimSpace(sort))
	if sort == "" {
		sort = "desc"
	}
	return fmt.Sprintf("users:%s:%s:%s:%d:%d", search, orderBy, sort, page, limit)
}
