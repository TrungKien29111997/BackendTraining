package repository

import "sever4/internal/models"

type UserRepository interface {
	Create(user models.User) error
	FindByUUID()
	Update()
	Delete()
	FillAll()
	FindByEmail(email string) (models.User, bool)
}
