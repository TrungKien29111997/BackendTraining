package routes

import (
	"sever4/internal/handler"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *handler.UserHandler
}

func NewUserRouter(handler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}

func (ur *UserRoutes) Register(r *gin.RouterGroup) {
	user := r.Group("/users")
	user.GET("", ur.handler.GetAllUser)
	user.POST("", ur.handler.CreateUser)
	user.GET("/:uuid", ur.handler.GetUserByUUID)
	user.PUT("/:uuid", ur.handler.UpdateUser)
	user.DELETE("/:uuid", ur.handler.DeleteUser)
}
