package service

import (
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
	us.repo.FillAll()
}
func (us *userService) CreateUser() {

}
func (us *userService) GetUserByUUID() {

}
func (us *userService) UpdateUser() {

}
func (us *userService) DeleteUser() {

}
