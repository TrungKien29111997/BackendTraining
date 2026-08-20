package v1handler

import (
	"net/http"
	v1dto "user-management-api/internal/dto/v1"
	v1service "user-management-api/internal/service/v1"
	"user-management-api/internal/utils"
	"user-management-api/internal/validation"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service v1service.AuthService
}

func NewAuthHandler(service v1service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}
func (uh *AuthHandler) CreateUser(ctx *gin.Context) {
	var input *v1dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	mapUserInput := input.MapCreateInputToModel()
	createUser, err := uh.service.CreateUser(ctx, mapUserInput)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := v1dto.MapUserToDTO(createUser)

	utils.ResponseSuccess(ctx, http.StatusCreated, userDTO)
}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	var input v1dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	_, accessToken, expiresIn, err := ah.service.Login(ctx, input.Email, input.Password)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	response := v1dto.LoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}
	//userDTO := v1dto.MapUserToDTO(authUser)
	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
func (ah *AuthHandler) Logout(ctx *gin.Context) {

}
