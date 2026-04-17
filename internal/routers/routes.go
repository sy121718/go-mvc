package routers

import (
	adminRouter "go-mvc/internal/module/backend/admin/router"
	r "go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

/*
路由注册 - SetupRoutes
===========================================
PHP 对比：
- Laravel: Route::group(['prefix' => 'api'], function() { ... })
- Go: 使用路由组（Group）组织路由

知识点：
1. 路由组 - 按功能分组（如 /api 前缀）
2. 中间件 - 全局中间件、路由组中间件
3. 路由方法 - GET/POST（本项目只用这两种）
*/

// SetupRoutes 注册所有路由
func SetupRoutes(router *gin.Engine) {
	/*
		健康检查接口
		==========================
		用途：检查服务是否正常运行
		路径：GET /health

		PHP 对比：
		Route::get('/health', function() {
			return response()->json(['status' => 'ok']);
		});
	*/
	router.GET("/health", func(c *gin.Context) {
		r.SuccessWithMessage(c, "msg_operation_success", gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api")
	{
		adminRouter.SetupAdminRoutes(api)
	}

	/*
		404 处理
		==========================
		当访问未定义的路由时返回 404

		使用统一响应格式：
		r.NotFound(c, "请求的资源不存在")
	*/
	router.NoRoute(func(c *gin.Context) {
		r.NotFound(c, "请求的资源不存在")
	})
}
