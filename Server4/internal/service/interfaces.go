package service

import "sever4/internal/models"

type UserService interface {
	GetAllUser(Search string, Page int, Limit int) ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetUserByUUID(uuid string) (models.User, error)
	UpdateUser(uuid string, user models.User) (models.User, error)
	DeleteUser(uuid string) error
}
