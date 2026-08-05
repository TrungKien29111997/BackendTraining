package repository

import "log"

type SQLUserRepository struct {
}

func NewSQLUserRepository() UserRepository {
	return &SQLUserRepository{}
}

func (u *SQLUserRepository) Create() {
	log.Println("Create")
}

func (u *SQLUserRepository) Find() {
	log.Println("Find")
}