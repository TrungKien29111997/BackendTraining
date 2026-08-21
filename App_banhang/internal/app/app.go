package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/db"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/routes"
	"user-management-api/internal/validation"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Module interface {
	Routes() routes.Route
}
type ModuleContext struct {
	DB    sqlc.Querier
	Redis *redis.Client
}

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}

func NewApplication(cfg *config.Config) *Application {

	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Validator init failed %v", err)
	}

	loadEnv()

	r := gin.Default()

	if err := db.InitDB(); err != nil {
		log.Fatalf("DB init failed %v", err)
	}

	redisClient := config.NewRedisClient()

	ctx := &ModuleContext{
		DB:    db.DB,
		Redis: redisClient,
	}
	cacheRedisService := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheRedisService)

	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService),
	}

	routes.RegisterRoutes(r, tokenService, getModulRoutes(modules)...)

	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
	}
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: a.router,
	}

	quit := make(chan os.Signal, 1)
	// sigcall.SIGINT => Ctrl + C
	// sigcall.SIGTERM => kill server
	// sigcall.SIGHUP => Reload server
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		log.Printf("Server is running %s \n", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server success shutdown")
	return nil
}

func getModulRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}
	return routeList
}

func loadEnv() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get working dir:", err)
	}
	evnPath := filepath.Join(cwd, ".env")

	err = godotenv.Load(evnPath)
	if err != nil {
		log.Println("No .env file found")
	} else {
		log.Println("Load success .env")
	}
}
