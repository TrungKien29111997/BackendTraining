package dto

import "sever4/internal/models"

type UserDTO struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Age    int    `json:"age"`
	Status string `json:"status"`
	Level  string `json:"level"`
}
type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Age      int    `json:"age" binding:"required"`
	Password string `json:"password" binding:"required"`
	Status   int    `json:"status" binding:"required"`
	Level    int    `json:"level" binding:"required"`
}
type UpdateUserInput struct {
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty"`
	Age      int    `json:"age" binding:"omitempty"`
	Password string `json:"password" binding:"omitempty"`
	Status   int    `json:"status" binding:"omitempty"`
	Level    int    `json:"level" binding:"omitempty"`
}

func MapUserToDTO(user models.User) *UserDTO {
	return &UserDTO{
		UUID:   user.UUID,
		Name:   user.Name,
		Email:  user.Email,
		Age:    user.Age,
		Status: mapStatusText(user.Status),
		Level:  mapLevelText(user.Level),
	}
}
func MapUsersToDTO(users []models.User) *[]UserDTO {
	dtos := make([]UserDTO, 0, len(users))
	for _, user := range users {
		dtos = append(dtos, *MapUserToDTO(user))
	}
	return &dtos
}
func (input *CreateUserInput) MapCreateInputToModel() models.User {
	return models.User{
		Name:     input.Name,
		Email:    input.Email,
		Age:      input.Age,
		Password: input.Password,
		Status:   input.Status,
		Level:    input.Level,
	}
}
func (input *UpdateUserInput) MapUpdateInputToModel() models.User {
	return models.User{
		Name:     input.Name,
		Email:    input.Email,
		Age:      input.Age,
		Password: input.Password,
		Status:   input.Status,
		Level:    input.Level,
	}
}

func mapStatusText(status int) string {
	switch status {
	case 1:
		return "Show"
	case 2:
		return "Hide"
	default:
		return "None"
	}
}

func mapLevelText(level int) string {
	switch level {
	case 1:
		return "Admin"
	case 2:
		return "User"
	default:
		return "None"
	}
}
