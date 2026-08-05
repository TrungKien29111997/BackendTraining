package repository

import (
	"hoc-gin/internal/models"
	"log"
)

type SQLUserRepository struct {
}

func NewSQLUserRepository() UserRepository {
	return &SQLUserRepository{}
}

func (u *SQLUserRepository) Create(user models.User) {
	log.Println("Create")
}

func (u *SQLUserRepository) Find(id int) {
	log.Println("Find")
}
