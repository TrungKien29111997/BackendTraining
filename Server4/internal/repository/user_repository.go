package repository

import (
	"log"
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
	log.Println("Get All Users")
}
func (ur *InMemoryUserRepository) Create() {

}
func (ur *InMemoryUserRepository) FindByUUID() {

}
func (ur *InMemoryUserRepository) Update() {

}
func (ur *InMemoryUserRepository) Delete() {

}
