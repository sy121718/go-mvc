package adminrouter

import "github.com/gin-gonic/gin"

// SetupAdminRoutes 注册 admin 模块路由
func SetupAdminRoutes(rg *gin.RouterGroup) {
	rg.Group("/admin")
}
