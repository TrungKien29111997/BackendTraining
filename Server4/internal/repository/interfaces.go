package repository

import "sever4/internal/models"

type UserRepository interface {
	Create(user models.User) error
	FindByUUID(uuid string) (models.User, bool)
	Update(uuid string, user models.User) error
	Delete(uuid string) error
	FillAll() ([]models.User, error)
	FindByEmail(email string) (models.User, bool)
}
