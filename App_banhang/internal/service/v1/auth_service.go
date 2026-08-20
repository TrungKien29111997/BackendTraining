package v1service

import (
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
	repo         repository.UserRepository
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService) *authService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
		repo:         repo,
	}
}
func (us *authService) CreateUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
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
	// if err := us.cache.Clear("users:*"); err != nil {
	// 	log.Printf("Failed to clear cache: %v", err)
	// }
	return user, nil
}
func (as *authService) Login(ctx *gin.Context, email, password string) (sqlc.User, string, int, error) {
	context := ctx.Request.Context()

	email = utils.NormalizeString(email)

	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		return sqlc.User{}, "", 0, utils.NewError("Failed to login", utils.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return sqlc.User{}, "", 0, utils.NewError("Failed to login", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return sqlc.User{}, "", 0, utils.WrapError(err, "Failed to generate access token", utils.ErrCodeInternal)
	}

	return user, accessToken, int(auth.AccessTokenTTL.Seconds()), nil
}
func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}
