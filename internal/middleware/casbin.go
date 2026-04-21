package middleware

import (
	"strconv"

	"go-mvc/pkg/casbin"
	"go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithMessage(c, 401, "未获取到用户信息")
			c.Abort()
			return
		}

		obj := c.Request.URL.Path
		act := c.Request.Method
		sub := strconv.FormatInt(userID.(int64), 10)

		enforcer := casbin.GetEnforcer()
		if enforcer == nil {
			response.ErrorWithMessage(c, 500, "权限系统未初始化")
			c.Abort()
			return
		}

		ok, err := enforcer.Enforce(sub, obj, act)
		if err != nil {
			response.ErrorWithMessage(c, 500, "权限验证失败")
			c.Abort()
			return
		}

		if !ok {
			response.ErrorWithMessage(c, 403, "无权限访问")
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(int64)
}

func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}
