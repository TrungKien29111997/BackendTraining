package v1service

import (
	"strings"
	"sync"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
	repo         repository.UserRepository
	cache        cache.RedisCacheService
}
type LoginAttempt struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu              sync.Mutex
	clients         = make(map[string]*LoginAttempt)
	LoginAttemptTTL = 5 * time.Minute
	MaxLoginAttempt = 5
)

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService, cache cache.RedisCacheService) *authService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
		repo:         repo,
		cache:        cache,
	}
}

func (as *authService) getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func (as *authService) getLoginAttempt(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(float32(MaxLoginAttempt)/float32(LoginAttemptTTL.Seconds())), MaxLoginAttempt) // 5 request/sec, brust 10
		newClient := &LoginAttempt{limiter, time.Now()}
		clients[ip] = newClient
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

func (as *authService) checkLoginAttempt(ip string) error {
	limiter := as.getLoginAttempt(ip)
	if !limiter.Allow() {
		return utils.NewError("Too many login attempt", utils.ErrCodeTooManyRequest)
	}
	return nil
}

func (as *authService) CleanupClients(ip string) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, ip)
}

func (as *authService) CreateUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()
	input.UserEmail = utils.NormalizeString(input.UserEmail)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "Failed to hash password", utils.ErrCodeInternal)
	}
	input.UserPassword = string(hashedPassword)
	user, err := as.repo.Create(context, input)
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
	ip := as.getClientIP(ctx)

	if err := as.checkLoginAttempt(ip); err != nil {
		return sqlc.User{}, "", "", 0, err
	}

	email = utils.NormalizeString(email)

	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		as.getLoginAttempt(ip)
		return sqlc.User{}, "", "", 0, utils.NewError("Failed to login", utils.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		as.getLoginAttempt(ip)
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
	as.CleanupClients(ip)
	return user, accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
func (as *authService) Logout(ctx *gin.Context, refreshToken string) error {

	//context := ctx.Request.Context()

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return utils.NewError("Unauthorized", utils.ErrCodeUnauthorized)
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	_, claims, err := as.tokenService.ParseToken(accessToken)
	if err != nil {
		return utils.NewError("Invalid access token", utils.ErrCodeUnauthorized)
	}
	if jti, ok := claims["jti"].(string); ok {
		expUnix, _ := claims["exp"].(float64)
		exp := time.Unix(int64(expUnix), 0)
		key := "blacklist:" + jti
		ttl := time.Until(exp)
		as.cache.Set(key, "revoked", ttl)
	}

	// parse refresh token
	token, err := as.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	_, err = uuid.Parse(token.UserUUID)
	if err != nil {
		return utils.NewError("Failed to parse user uuid", utils.ErrCodeUnauthorized)
	}
	// revoke old refresh token
	if err := as.tokenService.RevokeRefreshToken(refreshToken); err != nil {
		return err
	}

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
