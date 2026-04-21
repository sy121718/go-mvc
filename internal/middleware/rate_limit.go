package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	count   int
	resetAt time.Time
}

var (
	rateLimitMu    sync.Mutex
	rateLimitStore = map[string]rateLimitEntry{}
)

func RequestRateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limit <= 0 || window <= 0 {
			c.Next()
			return
		}

		key := buildRateLimitKey(c)
		now := time.Now()

		rateLimitMu.Lock()
		entry, exists := rateLimitStore[key]
		if !exists || !entry.resetAt.After(now) {
			entry = rateLimitEntry{
				count:   0,
				resetAt: now.Add(window),
			}
		}
		entry.count++
		rateLimitStore[key] = entry
		rateLimitMu.Unlock()

		if entry.count > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Response{
				Code:    http.StatusTooManyRequests,
				Message: "请求过于频繁",
			})
			return
		}

		c.Next()
	}
}

func buildRateLimitKey(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return "unknown"
	}
	path := strings.TrimSpace(c.FullPath())
	if path == "" {
		path = strings.TrimSpace(c.Request.URL.Path)
	}
	return c.ClientIP() + "|" + path
}

func resetRateLimitStore() {
	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()
	rateLimitStore = map[string]rateLimitEntry{}
}
