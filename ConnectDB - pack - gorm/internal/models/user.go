package models

type User struct {
	Id    int    `json:"user_id" gorm:"column:user_id;primary_key"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
