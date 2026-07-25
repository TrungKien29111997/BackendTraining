package repository

import (
	"fmt"
	"sever4/internal/models"
	"slices"
)

type InMemoryUserRepository struct {
	user []models.User
}

func NewUserRepository() UserRepository {
	return &InMemoryUserRepository{
		user: make([]models.User, 0),
	}
}
func (ur *InMemoryUserRepository) FillAll() ([]models.User, error) {
	return ur.user, nil
}
func (ur *InMemoryUserRepository) Create(user models.User) error {
	ur.user = append(ur.user, user)
	return nil
}
func (ur *InMemoryUserRepository) FindByUUID(uuid string) (models.User, bool) {
	for _, user := range ur.user {
		if user.UUID == uuid {
			return user, true
		}
	}
	return models.User{}, false
}
func (ur *InMemoryUserRepository) FindByEmail(email string) (models.User, bool) {
	for _, user := range ur.user {
		if user.Email == email {
			return user, true
		}
	}
	return models.User{}, false
}
func (ur *InMemoryUserRepository) Update(uuid string, user models.User) error {
	for i, u := range ur.user {
		if u.UUID == uuid {
			ur.user[i] = user
			return nil
		}
	}
	return fmt.Errorf("user not found")
}
func (ur *InMemoryUserRepository) Delete(uuid string) error {
	for i, u := range ur.user {
		if u.UUID == uuid {
			ur.user = slices.Delete(ur.user, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("user not found")
}
