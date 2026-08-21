package v1service

import (
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func (as *authService) Login(ctx *gin.Context, email, password string) (sqlc.User, string, string, int, error) {
	context := ctx.Request.Context()

	email = utils.NormalizeString(email)

	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		return sqlc.User{}, "", "", 0, utils.NewError("Failed to login", utils.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return sqlc.User{}, "", "", 0, utils.NewError("Failed to login", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return sqlc.User{}, "", "", 0, utils.WrapError(err, "Failed to generate access token", utils.ErrCodeInternal)
	}

	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return sqlc.User{}, "", "", 0, utils.WrapError(err, "Failed to generate access token", utils.ErrCodeInternal)
	}

	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return sqlc.User{}, "", "", 0, utils.WrapError(err, "Cannot store refresh token", utils.ErrCodeInternal)
	}

	return user, accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}
func (as *authService) RefreshToken(ctx *gin.Context, tokenStr string) (string, string, int, error) {

	context := ctx.Request.Context()

	// parse refresh token
	token, err := as.tokenService.ValidateRefreshToken(tokenStr)
	if err != nil {
		return "", "", 0, err
	}
	userUuid, err := uuid.Parse(token.UserUUID)
	if err != nil {
		return "", "", 0, utils.NewError("Failed to parse user uuid", utils.ErrCodeUnauthorized)
	}

	// get user
	user, err := as.userRepo.GetByUUID(context, userUuid)
	if err != nil {
		return "", "", 0, utils.NewError("Failed to get user", utils.ErrCodeUnauthorized)
	}

	//generate new access token
	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Failed to generate access token", utils.ErrCodeInternal)
	}

	//generate new refresh token
	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Failed to generate refresh token", utils.ErrCodeInternal)
	}

	// revoke old refresh token
	if err := as.tokenService.RevokeRefreshToken(tokenStr); err != nil {
		return "", "", 0, err
	}

	// save to redis cache
	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError(err, "Cannot store refresh token", utils.ErrCodeInternal)
	}
	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
