package service

import (
	"sever4/internal/models"
	"sever4/internal/repository"
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
func (us *userService) CreateUser(user models.User) models.User {
	if _, exist := us.repo.FindByEmail(user.Email); exist {
		return models.User{}
	}
}
func (us *userService) GetUserByUUID() {

}
func (us *userService) UpdateUser() {

}
func (us *userService) DeleteUser() {

}
