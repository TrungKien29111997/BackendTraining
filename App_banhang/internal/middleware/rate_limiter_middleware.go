package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"
	"user-management-api/internal/utils"

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

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]
	if !exists {
		rateLimiterSecStr := utils.GetEnv("RATE_LIMITER_REQUEST_SEC", "5")
		rateLimiterBurstStr := utils.GetEnv("RATE_LIMITER_BURST", "10")
		rateLimiterSec, err := strconv.Atoi(rateLimiterSecStr)
		if err != nil {
			panic("invalid RATE_LIMITER_REQUEST_SEC: " + err.Error())
		}
		rateLimiterBurst, err := strconv.Atoi(rateLimiterBurstStr)
		if err != nil {
			panic("invalid RATE_LIMITER_BURST: " + err.Error())
		}
		limiter := rate.NewLimiter(rate.Limit(rateLimiterSec), rateLimiterBurst) // 5 request/sec, brust 10
		newClient := &Client{limiter, time.Now()}
		clients[ip] = newClient
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// ab -n 20 -c 1 -H "X-API-Key:ab2a7a8a-d601-4bf7-b0e2-dd00e5459392" http://localhost:8080/api/v1/categories/golang
func RateLimiterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getClientIP(ctx)

		limiter := getRateLimiter(ip)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many request",
				"message": "Bạn đã gửi quá nhiêu request. Hãy thử lại sau",
			})

			return
		}

		ctx.Next()
	}
}
