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
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService) *authService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
	}
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
