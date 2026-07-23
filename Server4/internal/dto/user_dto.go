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
