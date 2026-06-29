package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
}
type TestBodyJson struct {
	Name string
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (u *ProductHandler) GetProductV1(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	c.JSON(http.StatusOK, gin.H{
		"message": "List product v1",
		"limit":   limit,
	})
}
func (u *ProductHandler) PostProductV1(c *gin.Context) {
	var body TestBodyJson
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Create product v1",
		"body":    body.Name,
	})
}
