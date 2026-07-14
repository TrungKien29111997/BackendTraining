package main

import (
	"sever4/internal/config"
	"sever4/internal/handler"
	"sever4/internal/repository"
	"sever4/internal/routes"
	"sever4/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()

	userRepo := repository.NewUserRepository()

	userService := service.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	userRoutes := routes.NewUserRouter(userHandler)

	r := gin.Default()

	routes.RegisterRoutes(r, userRoutes)

	if err := r.Run(cfg.ServerAddress); err != nil {
		panic(err)
	}
}
