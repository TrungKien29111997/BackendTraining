package repository

import (
	"sever4/internal/models"
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
func (ur *InMemoryUserRepository) Update() {

}
func (ur *InMemoryUserRepository) Delete() {

}
