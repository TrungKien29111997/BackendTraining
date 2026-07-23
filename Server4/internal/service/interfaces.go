package service

import "sever4/internal/models"

type UserService interface {
	GetAllUser() ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetUserByUUID(uuid string) (models.User, error)
	UpdateUser()
	DeleteUser()
}
