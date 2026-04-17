package middleware

import (
	"go-mvc/pkg/casbin"
	"go-mvc/pkg/enums"
	"go-mvc/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CasbinMiddleware Casbin 权限验证中间件
// 从 Context 获取用户 ID，自动验证当前请求路径的权限
// sub = 用户ID, obj = 请求路径, act = 请求方法(GET/POST)
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Context 获取用户 ID（由 JWT 中间件存入）
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithMessage(c, enums.ErrUnauthorized, "未获取到用户信息")
			return
		}

		// 2. 获取请求路径和方法
		obj := c.Request.URL.Path // 资源：请求路径
		act := c.Request.Method   // 操作：GET 或 POST

		// 3. 将用户 ID 转为字符串作为 sub
		sub := strconv.FormatInt(userID.(int64), 10)

		// 4. 调用 Casbin 验证权限
		enforcer := casbin.GetEnforcer()
		if enforcer == nil {
			response.ErrorWithMessage(c, enums.ErrSystemError, "权限系统未初始化")
			return
		}

		// Enforce(sub, obj, act) 返回 true 表示有权限
		ok, err := enforcer.Enforce(sub, obj, act)
		if err != nil {
			response.ErrorWithMessage(c, enums.ErrSystemError, "权限验证失败")
			return
		}

		if !ok {
			response.ErrorWithMessage(c, enums.ErrPermissionDenied, "无权限访问")
			return
		}

		c.Next()
	}
}

// GetUserID 从 Context 获取用户 ID 的辅助函数
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetUsername 从 Context 获取用户名的辅助函数
func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}
