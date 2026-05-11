// Package builtin 内置 Gin HTTP 中间件实现。
//
// 每个文件对应一个独立的功能型中间件，所有中间件函数返回 gin.HandlerFunc。
// 本包不读取外部配置，所有参数由调用方传入（通常是 middleware.Setup()）。
// 本包不导入 go-mvc/config，保持纯粹性。
package builtin

import (
	"strings"

	"go-mvc/pkg/auth"
	"go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT 认证中间件。
//
// 校验流程：
//  1. 从 Authorization 头中提取 Bearer token
//  2. 调用 auth.ParseToken() 解析 JWT，验证签名和过期时间
//  3. 解析成功后，将 user_id 和 username 写入 gin.Context
//
// 失败场景：
//   - 未携带 Authorization 头 → 返回 401 "未登录或登录已过期"
//   - Authorization 格式错误（非 Bearer） → 返回 401
//   - token 无效或已过期 → 返回 401
//
// 适用位置：需要登录认证的路由组或单路由。
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithMessage(c, 401, "未登录或登录已过期")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorWithMessage(c, 401, "未登录或登录已过期")
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			response.ErrorWithMessage(c, 401, "未登录或登录已过期")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}