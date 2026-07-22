package service

import (
	"sever4/internal/models"
	"sever4/internal/repository"
	"sever4/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}
func (us *userService) GetAllUser() {
}
func (us *userService) CreateUser(user models.User) (models.User, error) {
	if _, exist := us.repo.FindByEmail(user.Email); exist {
		return models.User{}, utils.NewError("user already exist", utils.ErrConflict)
	}

	user.UUID = uuid.New().String()
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, utils.WrapError(err, "failed to hash password", utils.ErrInternalServerError)
	}

	user.Password = string(hashPass)
	if err := us.repo.Create(user); err != nil {
		return models.User{}, utils.WrapError(err, "failed to create user", utils.ErrInternalServerError)
	}

	return user, nil
}
func (us *userService) GetUserByUUID() {

}
func (us *userService) UpdateUser() {

}
func (us *userService) DeleteUser() {

}
