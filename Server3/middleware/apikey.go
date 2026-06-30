package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleWare() gin.HandlerFunc {
	expectedAPIKey := os.Getenv("API_KEY")
	if expectedAPIKey == "" {
		expectedAPIKey = "default"
	}
	log.Println("Config", expectedAPIKey)
	return func(c *gin.Context) {
		apikey := c.GetHeader("X-API-KEY")
		log.Println("APIKEY", apikey)
		if apikey == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Non API_KEY"})
			return
		}
		if apikey != expectedAPIKey {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong API_KEY"})
			return
		}
		c.Next()
	}
}
