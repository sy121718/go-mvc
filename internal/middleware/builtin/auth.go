package builtin

import (
	"log"
	"strings"
	"time"

	"go-mvc/pkg/auth"
	"go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT 认证与自动续期中间件。
//
// 认证校验：
//  1. 从 Authorization 头中提取 Bearer token
//  2. 调用 auth.ParseToken() 解析 JWT，验证签名和过期时间
//  3. 解析成功后，将 user_id 和 username 写入 gin.Context
//
// 自动续期：
//   请求处理完成后，若 token 剩余有效期不足 1 小时，
//   自动生成新 token 并通过 X-New-Token 响应头返回前端。
//   前端 Axios 拦截器检测到此头后自动更新本地存储，实现无感续期。
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

		// 自动续期：token 剩余不足 1 小时，生成新 token 通过响应头返回
		if claims.ExpiresAt != nil && time.Until(claims.ExpiresAt.Time) <= time.Hour {
			newToken, _, _, err := auth.GenerateTokenPair(claims.UserID, claims.Username, false)
			if err != nil {
				log.Printf("JWT 自动续期失败: %v", err)
				return
			}
			c.Header("X-New-Token", newToken)
			c.Header("Access-Control-Expose-Headers", "X-New-Token")
		}
	}
}