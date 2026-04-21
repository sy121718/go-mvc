package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// DefaultOptions 默认 HTTP 中间件链的配置项。
type DefaultOptions struct {
	RequestBodyLimit int64
	UploadBodyLimit  int64
	RateLimitEnabled bool
	RateLimitLimit   int
	RateLimitWindow  time.Duration
	LogCapture       bool
}

// UseDefaultMiddlewares 按统一顺序挂载框架默认中间件。
//
// 默认顺序：
// 1. Recovery
// 2. SecurityHeaders
// 3. RequestBodyLimit
// 4. RequestRateLimit
// 5. RequestLogCapture
func UseDefaultMiddlewares(engine *gin.Engine, options DefaultOptions) {
	if engine == nil {
		return
	}

	engine.Use(gin.Recovery())
	engine.Use(SecurityHeadersMiddleware())
	engine.Use(RequestBodyLimitMiddleware(options.RequestBodyLimit, options.UploadBodyLimit))
	if options.RateLimitEnabled {
		engine.Use(RequestRateLimitMiddleware(options.RateLimitLimit, options.RateLimitWindow))
	}
	engine.Use(RequestLogCaptureMiddleware(options.LogCapture))
}
