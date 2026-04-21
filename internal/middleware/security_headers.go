package middleware

import "github.com/gin-gonic/gin"

// SecurityHeadersMiddleware 为所有 HTTP 响应追加基础安全响应头。
//
// 目的：
// - 降低 MIME 嗅探风险
// - 禁止页面被 iframe 嵌套
// - 提供最基础的内容安全策略
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'; base-uri 'self'")
		c.Next()
	}
}
