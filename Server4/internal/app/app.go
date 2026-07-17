package app

import (
	"sever4/internal/config"
	"sever4/internal/routes"

	"github.com/gin-gonic/gin"
)

type Application struct {
	config *config.Config
	router *gin.Engine
}

type Module interface {
	Routes() routes.Route
}

func NewApplication(cfg *config.Config) *Application {
	r := gin.Default()

	modules := []Module{
		NewUserModule(),
	}

	routes.RegisterRoutes(r, getModuleRoutes(modules)...)
	return &Application{
		config: cfg,
		router: r,
	}
}

func (a *Application) Run() error {
	return a.router.Run(a.config.ServerAddress)
}

func getModuleRoutes(m []Module) []routes.Route {
	rotesList := make([]routes.Route, len(m))
	for i, module := range m {
		rotesList[i] = module.Routes()
	}
	return rotesList
}
