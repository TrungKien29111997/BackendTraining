package repository

import "sever4/internal/models"

type InMemoryUserRepository struct {
	user []models.User
}

func NewUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		user: make([]models.User, 0),
	}
}