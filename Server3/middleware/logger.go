package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func LoggerMiddleware() gin.HandlerFunc {
	logPath := "logs/http.log"
	if err := os.MkdirAll(filepath.Dir(logPath), os.ModePerm); err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	logger := zerolog.New(logFile).With().Timestamp().Logger()
	return func(c *gin.Context) {
		start := time.Now()
		contentType := c.GetHeader("Content-Type")
		requestBody := make(map[string]any)
		if strings.HasPrefix(contentType, "multipart/form-data") {
			log.Panicln("multipart/form-data")
		} else {
			// application/json
			// application/x-www-form-urlencoded
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Error().Str("error", err.Error()).Msg("read body error")
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			if strings.HasSuffix(contentType, "application/json") {
				json.Unmarshal(bodyBytes, &requestBody)
			} else {

			}
		}
		c.Next()
		duration := time.Since(start)

		logEvent := logger.Info()

		logEvent.Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Str("referer", c.Request.Referer()).
			Str("protocol", c.Request.Proto).
			Str("host", c.Request.Host).
			Str("content_type", c.ContentType()).
			Interface("headers", c.Request.Header).
			Int("status", c.Writer.Status()).
			Str("duration", duration.String()).
			Interface("body", requestBody).
			Msg("request completed")
	}
}
