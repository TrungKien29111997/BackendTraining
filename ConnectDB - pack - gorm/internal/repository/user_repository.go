package repository

import (
	"hoc-gin/internal/models"
	"log"

	"gorm.io/gorm"
)

type SQLUserRepository struct {
	db *gorm.DB
}

func NewSQLUserRepository(db *gorm.DB) UserRepository {
	return &SQLUserRepository{
		db: db,
	}
}

func (u *SQLUserRepository) Create(user *models.User) error {
	if err := u.db.Create(user).Error; err != nil {
		return err
	}
	log.Println("Create")
	return nil
}

func (u *SQLUserRepository) Find(id int, user *models.User) error {
	if err := u.db.First(user, id).Error; err != nil {
		return err
	}
	log.Println("Find")
	return nil
}
