package userrouter

import "github.com/gin-gonic/gin"

// SetupUserRoutes 注册 user 模块路由
func SetupUserRoutes(rg *gin.RouterGroup) {
	rg.Group("/user")
}
