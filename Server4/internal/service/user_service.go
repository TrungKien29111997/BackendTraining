package service

import (
	"sever4/internal/models"
	"sever4/internal/repository"
	"sever4/internal/utils"
	"strings"

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
func (us *userService) GetAllUser(Search string, Page int, Limit int) ([]models.User, error) {
	users, err := us.repo.FillAll()
	if err != nil {
		return nil, utils.WrapError(err, "failed to get all user", utils.ErrInternalServerError)
	}

	var filteredUsers []models.User
	if Search != "" {
		Search = strings.ToLower(Search)
		for _, user := range users {
			name := strings.ToLower(user.Name)
			email := strings.ToLower(user.Email)
			if strings.Contains(name, Search) || strings.Contains(email, Search) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	} else {
		filteredUsers = users
	}
	start := (Page - 1) * Limit
	if start > len(filteredUsers) {
		return []models.User{}, nil
	}
	end := Page + Limit
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}
	return filteredUsers[start:end], nil
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
func (us *userService) GetUserByUUID(uuid string) (models.User, error) {
	user, exist := us.repo.FindByUUID(uuid)
	if !exist {
		return models.User{}, utils.NewError("user not found", utils.ErrNotFound)
	}
	return user, nil
}
func (us *userService) UpdateUser(uuid string, updateUser models.User) (models.User, error) {

	if u, exist := us.repo.FindByEmail(updateUser.Email); exist && u.UUID != uuid {
		return models.User{}, utils.NewError("user already exist", utils.ErrConflict)
	}
	currentUser, found := us.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, utils.NewError("user not found", utils.ErrNotFound)
	}

	currentUser.Name = updateUser.Name
	currentUser.Email = updateUser.Email
	currentUser.Age = updateUser.Age
	currentUser.Status = updateUser.Status
	currentUser.Level = updateUser.Level

	if updateUser.Password != "" {
		hashPass, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, utils.WrapError(err, "failed to hash password", utils.ErrInternalServerError)
		}
		currentUser.Password = string(hashPass)
	}
	if err := us.repo.Update(uuid, currentUser); err != nil {
		return models.User{}, utils.WrapError(err, "failed to update user", utils.ErrInternalServerError)
	}
	return currentUser, nil
}
func (us *userService) DeleteUser(uuid string) error {
	if err := us.repo.Delete(uuid); err != nil {
		return utils.WrapError(err, "failed to delete user", utils.ErrInternalServerError)
	}
	return nil
}
