package repository

import "hoc-gin/internal/models"

type UserRepository interface {
	Create(user models.User)
	Find(id int)
}
