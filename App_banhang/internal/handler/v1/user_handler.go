package v1handler

import (
	"net/http"
	v1dto "user-management-api/internal/dto/v1"
	v1service "user-management-api/internal/service/v1"
	"user-management-api/internal/utils"
	"user-management-api/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service v1service.UserService
}

func NewUserHandler(service v1service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUsers(ctx *gin.Context) {
	var params v1dto.GetUsersParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	users, totalUsers, err := uh.service.GetAllUsers(ctx, params.Search, params.OrderBy, params.Sort, params.Page, params.Limit)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTOs := v1dto.MapUsersToDTO(users)
	paginationRespone := utils.NewPaginationRespone(userDTOs, params.Page, params.Limit, totalUsers)
	utils.ResponseSuccess(ctx, http.StatusOK, paginationRespone)
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
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

func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "")
}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	var input v1dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapUpdateInputToModel(userUuid)
	updatedUser, err := uh.service.UpdateUser(ctx, user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUserToDTO(updatedUser)
	utils.ResponseSuccess(ctx, http.StatusOK, userDTO)
}

func (uh *UserHandler) SoftDeleteUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	deletedUser, err := uh.service.SoftDeleteUser(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUserToDTO(deletedUser)
	utils.ResponseSuccess(ctx, http.StatusOK, userDTO)
}
func (uh *UserHandler) HardDeleteUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	_, err = uh.service.HardDeleteUser(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusOK)
}
func (uh *UserHandler) RestoreUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	restoreUser, err := uh.service.RestoreUser(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUserToDTO(restoreUser)
	utils.ResponseSuccess(ctx, http.StatusOK, userDTO)
}
