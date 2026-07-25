package handler

import (
	"net/http"
	"sever4/internal/dto"
	"sever4/internal/service"
	"sever4/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}
type GetUserByUUIDParam struct {
	Uiid string `uri:"uuid" binding:"uuid"`
}

type GetUserParams struct {
	Search string `form:"search" binding:"omitempty,min=3,max=50"`
	Page   int    `form:"page" binding:"omitempty,gte=1,lte=100"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}
func (uh *UserHandler) GetAllUser(c *gin.Context) {
	var params GetUserParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}
	users, err := uh.service.GetAllUser(params.Search, params.Page, params.Limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, "err get users")
		return
	}
	userDTOs := dto.MapUsersToDTO(users)
	utils.ReponseSuccess(c, http.StatusOK, &userDTOs)
}
func (uh *UserHandler) CreateUser(c *gin.Context) {
	var input dto.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := input.MapCreateInputToModel()
	createdUser, err := uh.service.CreateUser(user)
	if err != nil {
		utils.ResponeError(c, err)
		return
	}

	userDTO := dto.MapUserToDTO(createdUser)
	utils.ReponseSuccess(c, http.StatusCreated, &userDTO)
}

func (uh *UserHandler) GetUserByUUID(c *gin.Context) {
	var param GetUserByUUIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, "err get uuid user")
		return
	}

	user, err := uh.service.GetUserByUUID(param.Uiid)
	if err != nil {
		utils.ResponeError(c, err)
		return
	}
	userDTO := dto.MapUserToDTO(user)
	utils.ReponseSuccess(c, http.StatusOK, &userDTO)
}
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	var param GetUserByUUIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, "err get uuid user")
		return
	}
	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := input.MapUpdateInputToModel()
	updatedUser, err := uh.service.UpdateUser(param.Uiid, user)
	if err != nil {
		utils.ResponeError(c, err)
		return
	}
	userDTO := dto.MapUserToDTO(updatedUser)
	utils.ReponseSuccess(c, http.StatusOK, &userDTO)
}
func (uh *UserHandler) DeleteUser(c *gin.Context) {
	var param GetUserByUUIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, "err get uuid user")
		return
	}
	if err := uh.service.DeleteUser(param.Uiid); err != nil {
		utils.ResponeError(c, err)
		return
	}
	utils.ReponseStatusCode(c, http.StatusNoContent)
}
