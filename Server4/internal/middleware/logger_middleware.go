package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

func LoggerMiddleware() gin.HandlerFunc {
	logPath := "../../internal/logs/http.log"

	logger := zerolog.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     5,    //days
		Compress:   true, // disabled by default
		LocalTime:  true,
	}).With().Timestamp().Logger()

	return func(c *gin.Context) {
		start := time.Now()
		contentType := c.GetHeader("Content-Type")
		requestBody := make(map[string]any)
		var formfiles []map[string]any
		// multipart/form-data
		if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := c.Request.ParseMultipartForm(32 << 20); err != nil && c.Request.MultipartForm != nil {
				// for value
				for key, vals := range c.Request.MultipartForm.Value {
					requestBody[key] = vals
				}
			}

			// for file
			for field, files := range c.Request.MultipartForm.File {
				for _, file := range files {
					formfiles = append(formfiles, map[string]any{
						"filename":     file.Filename,
						"size":         file.Size,
						"field":        field,
						"Content-Type": file.Header.Get("Content-Type"),
					})
				}
			}
			requestBody["form_files"] = formfiles
		} else {

			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Error().Str("error", err.Error()).Msg("read body error")
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// application/json
			if strings.HasSuffix(contentType, "application/json") {
				json.Unmarshal(bodyBytes, &requestBody)
			} else {
				// application/x-www-form-urlencoded
				value, _ := url.ParseQuery(string(bodyBytes))
				for key, vals := range value {
					requestBody[key] = vals
				}
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
