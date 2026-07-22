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
func (ur *InMemoryUserRepository) FillAll() {

}
func (ur *InMemoryUserRepository) Create(user models.User) error {
	ur.user = append(ur.user, user)
	return nil
}
func (ur *InMemoryUserRepository) FindByUUID() {

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
