package routers

import "sever4/internal/handler"

type UserRoutes struct {
	handler *handler.UserHandler
}

func NewUserRouter(handler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}
