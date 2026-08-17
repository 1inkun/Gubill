package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter() gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu        sync.Mutex
		clients   = make(map[string]*client)
		startOnce sync.Once
	)

	// 定期清理长时间未访问的客户端，避免 map 无限增长
	startOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				mu.Lock()
				for ip, cl := range clients {
					if time.Since(cl.lastSeen) > 10*time.Minute {
						delete(clients, ip)
					}
				}
				mu.Unlock()
			}
		}()
	})

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		if _, exists := clients[ip]; !exists {
			// Allow 10 requests per second with a burst of 20
			clients[ip] = &client{limiter: rate.NewLimiter(10, 20), lastSeen: time.Now()}
		}
		cl := clients[ip]
		cl.lastSeen = time.Now()
		mu.Unlock()

		if !cl.limiter.Allow() {
			resp := models.NewResponse(429, "fail", "请求过于频繁", nil)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, resp)
			return
		}

		c.Next()
	}
}
