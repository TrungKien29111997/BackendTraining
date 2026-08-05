package repository

import (
	"database/sql"
	"fmt"
	"hoc-gin/internal/models"
	"log"
)

type SQLUserRepository struct {
	db *sql.DB
}

func NewSQLUserRepository(DB *sql.DB) UserRepository {
	return &SQLUserRepository{
		db: DB,
	}
}

func (u *SQLUserRepository) Create(user models.User) error {
	row := u.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email)
	err := row.Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}
	log.Println("Create")
	return nil
}

func (u *SQLUserRepository) Find(id int) {
	log.Println("Find")
}
