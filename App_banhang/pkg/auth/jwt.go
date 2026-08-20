package auth

import (
	"encoding/json"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
}

// type Claims struct {
// 	jwt.RegisteredClaims
// }

type EncryptedPayload struct {
	UserUUID  string `json:"user_uuid"`
	UserEmail string `json:"email"`
	Role      int32  `json:"role"`
}

var (
	jwtSecret     = []byte(utils.GetEnv("JWT_SECRET", "secret"))
	jwtEncryptKey = []byte(utils.GetEnv("JWT_ENCRYPT_KEY", "secret"))
)

const (
	AccessTokenTTL = 15 * time.Minute
)

func NewJWTService() TokenService {
	return &JWTService{}
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

func (js *JWTService) GenerateRefreshToken() string {
	return ""
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
