package app

import (
	"sever4/internal/handler"
	"sever4/internal/repository"
	"sever4/internal/routes"
	"sever4/internal/service"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewUserRepository()

	userService := service.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	userRoutes := routes.NewUserRouter(userHandler)

	return &UserModule{routes: userRoutes}
}

func (u *UserModule) Routes() routes.Route {
	return u.routes
}
