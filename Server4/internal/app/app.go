package app

import (
	"log"
	"sever4/internal/config"
	"sever4/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}

type Module interface {
	Routes() routes.Route
}

func NewApplication(cfg *config.Config) *Application {
	r := gin.Default()
	loadEnv()
	modules := []Module{
		NewUserModule(),
	}

	routes.RegisterRoutes(r, getModuleRoutes(modules)...)
	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
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

func loadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("no env file")
	}
}
