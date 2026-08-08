package routes

import (
	"user-management-api/internal/middleware"
	"user-management-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {

	httpLogger := newLoggerWhitePath("../../internal/logs/http.log", "info")

	recoveryLogger := newLoggerWhitePath("../../internal/logs/recovery.log", "warning")

	r.Use(
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
		middleware.RateLimiterMiddleware(),
	)

	v1api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1api)
	}
}

func newLoggerWhitePath(path string, level string) *zerolog.Logger {
	config := logger.LoggerConfig{
		Level:      level,
		Filename:   path,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     5, //days
		Compress:   true,
		IsDev:      "Dev", // utils.GetEnv("APP_EVN", "Dev"),
	}
	return logger.NewLogger(config)

}
