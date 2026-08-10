package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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

// hey -n 25 -c 1 -H "X-API-KEY: 1234" http://localhost:8080/api/v1/users
func RateLimiterMiddleware(rateLimiterLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getClientIP(ctx)

		limiter := getRateLimiter(ip)

		if !limiter.Allow() {
			if shouldLogRateLimit(ip) {
				rateLimiterLogger.Warn().
					Str("method", ctx.Request.Method).
					Str("path", ctx.Request.URL.Path).
					Str("query", ctx.Request.URL.RawQuery).
					Str("client_ip", ctx.ClientIP()).
					Str("user_agent", ctx.Request.UserAgent()).
					Str("referer", ctx.Request.Referer()).
					Str("protocol", ctx.Request.Proto).
					Str("host", ctx.Request.Host).
					Str("remote_addr", ctx.Request.RemoteAddr).
					Interface("headers", ctx.Request.Header).
					Str("request_uri", ctx.Request.RequestURI).
					Msg("Rate limiter Log")
			}
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many request",
				"message": "Bạn đã gửi quá nhiêu request. Hãy thử lại sau",
			})
			return
		}

		ctx.Next()
	}
}

var rateLimiterLogCache = sync.Map{}

const rateLimiterLogCacheDuration = 10 * time.Second

func shouldLogRateLimit(ip string) bool {
	now := time.Now()
	if v, ok := rateLimiterLogCache.Load(ip); ok {
		lastLogTime := v.(time.Time)
		if now.Sub(lastLogTime) < rateLimiterLogCacheDuration {
			return false
		}
	}
	rateLimiterLogCache.Store(ip, now)
	return true
}
