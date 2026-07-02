package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func GetClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "" {
		ip = c.Request.RemoteAddr
	}
	return ip
}

func GetRateLimmiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	client, exits := clients[ip]
	if !exits {
		limiter := rate.NewLimiter(5, 10) // 5 req/s, brust 10 (token)
		newClient := &Client{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		clients[ip] = newClient
		return limiter
	}
	client.lastSeen = time.Now()
	return client.limiter
}

func CleanUpClients() {
	for {
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
		time.Sleep(30 * time.Second)
	}
}

// hey -n 20 -c 1 -H "X-API-KEY: 1234" http://localhost:8080/api/v1/products
func RateLimitMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := GetClientIP(c)
		limiter := GetRateLimmiter(ip)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}
		c.Next()
	}
}
