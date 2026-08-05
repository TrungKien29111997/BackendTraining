package repository

import (
	"database/sql"
	"fmt"
	"hoc-gin/internal/models"
)

type SQLUserRepository struct {
	db *sql.DB
}

func NewSQLUserRepository(DB *sql.DB) UserRepository {
	return &SQLUserRepository{
		db: DB,
	}
}

func (u *SQLUserRepository) Create(user *models.User) error {
	row := u.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING user_id", user.Name, user.Email)
	err := row.Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}
	return nil
}

func (u *SQLUserRepository) Find(id int, user *models.User) error {
	row := u.db.QueryRow("SELECT * FROM users WHERE user_id = $1", id)
	err := row.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}
	return nil
}
