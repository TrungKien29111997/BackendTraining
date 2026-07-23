package handler

import (
	"net/http"
	"sever4/internal/dto"
	"sever4/internal/models"
	"sever4/internal/service"
	"sever4/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}
func (uh *UserHandler) GetAllUser(c *gin.Context) {
	users, err := uh.service.GetAllUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, "err get users")
		return
	}
	userDTOs := dto.MapUsersToDTO(users)
	utils.ReponseSuccess(c, http.StatusOK, &userDTOs)
}
func (uh *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdUser, err := uh.service.CreateUser(user)
	if err != nil {
		utils.ResponeError(c, err)
		return
	}

	userDTO := dto.MapUserToDTO(createdUser)
	utils.ReponseSuccess(c, http.StatusCreated, &userDTO)
}

type GetUserByUUIDParam struct {
	Uiid string `uri:"uuid" binding:"uuid"`
}

func (uh *UserHandler) GetUserByUUID(c *gin.Context) {
	var param GetUserByUUIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, "err get uuid user")
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

}
func (uh *UserHandler) DeleteUser(c *gin.Context) {

}
