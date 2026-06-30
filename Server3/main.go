package main

import (
	"Server3/middleware"
	v1handler "Server3/v1/handler"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("no env file")
	}

	r := gin.Default()

	r.Use(middleware.ApiKeyMiddleWare())

	v1 := r.Group("/api/v1")
	{
		newUser := v1handler.NewUserHandler()
		v1.GET("/users/:id", newUser.GetUserV1)
		v1.POST("/users", newUser.PostUserV1)

		newProduct := v1handler.NewProductHandler()
		v1.GET("/products", newProduct.GetProductV1)
		v1.POST("/products", newProduct.PostProductV1)
	}
	r.Run(":8080")
}
