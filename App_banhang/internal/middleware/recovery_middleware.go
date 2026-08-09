package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RecoveryMiddleware(recoveryLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				stactAt := ExtractFirstAppStackLine(stack)
				recoveryLogger.Error().
					Str("path", ctx.Request.URL.Path).
					Str("method", ctx.Request.Method).
					Str("client_ip", ctx.ClientIP()).
					Str("panic", fmt.Sprintf("%v", err)).
					//Str("stack", string(stack)).
					Str("stack_at", stactAt).
					Msg("Panic occurred")

				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    "Internal server error",
					"message": "Please try again later.",
				})
			}
		}()
		ctx.Next()
	}
}

func ExtractFirstAppStackLine(stack []byte) string {
	lines := bytes.Split(stack, []byte("\n"))
	for _, line := range lines {
		if bytes.Contains(line, []byte(".go")) &&
			!bytes.Contains(line, []byte("/runtime/")) &&
			!bytes.Contains(line, []byte("/debug/")) &&
			!bytes.Contains(line, []byte("recovery_middleware.go")) {
			cleanline := strings.TrimSpace(string(line))
			return cleanline
		}
	}
	return ""
}
