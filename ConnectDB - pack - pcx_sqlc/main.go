package main

import (
	"hoc-gin/internal/db"
	"hoc-gin/internal/handlers"
	"hoc-gin/internal/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	if err := db.InitDB(); err != nil {
		log.Fatal("unable to connect to db")
	}
	log.Println(db.DB)

	r := gin.Default()
	userRepo := repository.NewSQLUserRepository()
	user := handlers.NewUserHandler(userRepo)
	r.GET("/api/v1/users/:id", user.GetUserByUuid)
	r.POST("/api/v1/users", user.CreateUser)

	r.Run(":8080")
}
