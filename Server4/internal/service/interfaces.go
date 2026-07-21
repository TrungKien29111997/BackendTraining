package service

import "sever4/internal/models"

type UserService interface {
	GetAllUser()
	CreateUser(user models.User) models.User
	GetUserByUUID()
	UpdateUser()
	DeleteUser()
}
