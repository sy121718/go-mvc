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

// RequestRateLimitMiddleware 按“客户端 IP + 路由路径”进行固定窗口限流。
//
// 规则：
// - 同一 IP 对同一路径在一个窗口内最多访问 limit 次
// - 超过限制后返回 429
// - 适合作为基础框架默认限流，不依赖额外组件
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
				Code:    "ErrRateLimited",
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
