package handler

import (
	"sever4/internal/service"

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
	uh.service.GetAllUser()
}
func (uh *UserHandler) CreateUser(c *gin.Context) {

}
func (uh *UserHandler) GetUserByUUID(c *gin.Context) {

}
func (uh *UserHandler) UpdateUser(c *gin.Context) {

}
func (uh *UserHandler) DeleteUser(c *gin.Context) {

}
