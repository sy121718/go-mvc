package route

import (
	"github.com/gin-gonic/gin"
	adminhandle "go-mvc/internal/module/backend/admin/handle"
)

/*
路由注册 - SetupAdminRoutes

知识点：
1. 路由分组 - 按模块组织路由
2. 路由方法 - 只用 GET/POST
3. 路由路径 - 不需要写完整路径，由主路由控制前缀
*/

// SetupAdminRoutes 注册 admin 模块的路由
func SetupAdminRoutes(rg *gin.RouterGroup) {
	/*
		Admin 路由组
		==========================
		路径：/api/admin/*

		rg 是从主路由传入的路由组，已经有 /api 前缀
		这里再分组为 /admin，最终路径是 /api/admin/xxx
	*/
	admin := rg.Group("/admin")
	{
		// 测试接口 - GET 请求
		// 完整路径：GET /api/admin/test
		admin.GET("/test", adminhandle.Test)
		admin.GET("/test2", adminhandle.Test2)

	}
}
