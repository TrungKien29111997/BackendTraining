package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"
	"user-management-api/pkg/cache"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	cache cache.RedisCacheService
}

type EncryptedPayload struct {
	UserUUID  string `json:"user_uuid"`
	UserEmail string `json:"email"`
	Role      int32  `json:"role"`
}
type RefreshToken struct {
	Token     string    `json:"token"`
	UserUUID  string    `json:"user_uuid"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

var (
	jwtSecret     = []byte(utils.GetEnv("JWT_SECRET", "secret"))
	jwtEncryptKey = []byte(utils.GetEnv("JWT_ENCRYPT_KEY", "secret"))
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

func NewJWTService(cache cache.RedisCacheService) TokenService {
	return &JWTService{
		cache: cache,
	}
}

func (js *JWTService) GenerateAccessToken(user sqlc.User) (string, error) {
	payload := EncryptedPayload{
		UserUUID:  user.UserUuid.String(),
		UserEmail: user.UserEmail,
		Role:      user.UserLevel,
	}

	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encrypted, err := utils.EncryptAES(rawData, jwtEncryptKey)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"data": encrypted,
		"jti":  uuid.NewString(),
		"exp":  time.Now().Add(AccessTokenTTL).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "prj-banhang",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (js *JWTService) GenerateRefreshToken(user sqlc.User) (RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return RefreshToken{}, err
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	refeshToken := RefreshToken{
		Token:     token,
		UserUUID:  user.UserUuid.String(),
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked:   false,
	}
	return refeshToken, nil
}

func (js *JWTService) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}
	return token, claims, nil
}

func (js *JWTService) DecryptAccessTokenPayload(tokenString string) (*EncryptedPayload, error) {
	_, claims, err := js.ParseToken(tokenString)
	if err != nil {
		return nil, utils.WrapError(err, "Cannot parse token", utils.ErrCodeInternal)
	}
	raw, ok := claims["data"].(string)
	if !ok {
		return nil, utils.NewError("Encoded data not found", utils.ErrCodeUnauthorized)
	}
	decrypted, err := utils.DecryptAES(raw, jwtEncryptKey)
	if err != nil {
		return nil, utils.WrapError(err, "Cannot decrypt data", utils.ErrCodeInternal)
	}
	var payload *EncryptedPayload
	if err := json.Unmarshal(decrypted, &payload); err != nil {
		return nil, utils.WrapError(err, "Cannot unmarshal data", utils.ErrCodeInternal)
	}
	return payload, nil
}

func (js *JWTService) StoreRefreshToken(token RefreshToken) error {
	cacheKey := "refresh_token:" + token.Token
	return js.cache.Set(cacheKey, token, RefreshTokenTTL)
}

func (js *JWTService) ValidateRefreshToken(token string) (RefreshToken, error) {
	cacheKey := "refresh_token:" + token
	var refreshToken RefreshToken
	if err := js.cache.Get(cacheKey, &refreshToken); err != nil {
		return RefreshToken{}, utils.NewError("Refresh token not found", utils.ErrCodeInternal)
	}
	if refreshToken.Revoked {
		return RefreshToken{}, utils.NewError("Refresh token revoked", utils.ErrCodeUnauthorized)
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return RefreshToken{}, utils.NewError("Refresh token expired", utils.ErrCodeUnauthorized)
	}
	return refreshToken, nil
}

func (js *JWTService) RevokeRefreshToken(token string) error {
	cacheKey := "refresh_token:" + token
	var refreshToken RefreshToken
	if err := js.cache.Get(cacheKey, &refreshToken); err != nil {
		return utils.NewError("Refresh token not found", utils.ErrCodeInternal)
	}
	refreshToken.Revoked = true
	return js.cache.Set(cacheKey, refreshToken, time.Until(refreshToken.ExpiresAt))
}
