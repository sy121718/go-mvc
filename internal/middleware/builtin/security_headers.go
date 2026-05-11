package builtin

import "github.com/gin-gonic/gin"

// SecurityHeadersMiddleware 安全响应头中间件。
//
// 追加的响应头：
//   - X-Content-Type-Options: nosniff
//     禁止浏览器进行 MIME 类型嗅探，降低脚本注入风险
//   - X-Frame-Options: DENY
//     禁止页面被嵌套到 iframe 中，防止点击劫持
//   - Content-Security-Policy: default-src 'self'; frame-ancestors 'none'; base-uri 'self'
//     限制资源加载来源，仅允许同源资源
//
// 适用位置：全局 engine.Use()，确保所有响应都携带安全头。
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'; base-uri 'self'")
		c.Next()
	}
}