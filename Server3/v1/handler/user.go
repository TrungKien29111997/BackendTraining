package v1handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) GetUserV1(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user info: " + idStr,
	})
}
func (u *UserHandler) PostUserV1(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Create user v1",
	})
}
