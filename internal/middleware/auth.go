package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go-mvc/pkg/enums"
	"go-mvc/pkg/auth"
	"go-mvc/pkg/response"
)

// JWTAuthMiddleware JWT 认证中间件
// 从请求 Header 中提取 Token，验证并解析，将用户信息存入 Context
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithMessage(c, enums.ErrUnauthorized, "缺少认证信息")
			return
		}

		// 2. 解析 Bearer Token 格式
		// 格式：Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorWithMessage(c, enums.ErrUnauthorized, "认证格式错误")
			return
		}

		tokenString := parts[1]

		// 3. 解析 Token
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			response.ErrorWithMessage(c, enums.ErrInvalidToken, "Token 无效或已过期")
			return
		}

		// 4. 将用户信息存入 Context，供后续中间件和处理器使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
