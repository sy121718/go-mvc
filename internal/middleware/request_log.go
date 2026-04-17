package middleware

import (
	"errors"
	"fmt"
	"go-mvc/pkg/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogCaptureMiddleware 记录全局 HTTP 请求日志。
//
// enabled 表示是否开启捕获。
// 目录固定写入 public/logs/http/<yyyy-mm-dd>.log。
func RequestLogCaptureMiddleware(enabled bool) gin.HandlerFunc {
	if !enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	const scene = "http"
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)

		fields := map[string]interface{}{
			"method":     method,
			"path":       path,
			"status":     status,
			"client_ip":  c.ClientIP(),
			"latency_ms": latency.Milliseconds(),
		}
		if rawQuery != "" {
			fields["query"] = rawQuery
		}

		entry := logger.Scene(scene).WithFields(fields)
		if len(c.Errors) > 0 {
			errParts := make([]string, 0, len(c.Errors))
			for _, item := range c.Errors {
				errParts = append(errParts, item.Error())
			}
			joined := strings.Join(errParts, "; ")
			entry.With("errors", errParts).Error(errors.New(joined), "http request error")
			return
		}

		if status >= 500 {
			entry.Error(fmt.Errorf("http status %d", status), "http request failed")
			return
		}

		if status >= 400 {
			entry.Warn("http request warning")
			return
		}

		entry.Info("http request")
	}
}
