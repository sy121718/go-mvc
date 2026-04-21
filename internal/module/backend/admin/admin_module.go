package adminmodule

import (
	adminrouter "go-mvc/internal/module/backend/admin/router"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 通过模块统一入口注册 admin 模块路由。
//
// 说明：
// - config 层只依赖模块入口，不再直接依赖模块内部 router 子包
// - 模块内部如何组织 router/handle/service 由模块自己决定
func RegisterRoutes(rg *gin.RouterGroup) {
	adminrouter.SetupAdminRoutes(rg)
}
