package adminrouter

import (
	"go-mvc/internal/middleware/builtin"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(rg *gin.RouterGroup) {
	if rg == nil {
		return
	}

	handle, err := newAdminHandle()
	if err != nil {
		return
	}
	//绑定路由组
	admin := rg.Group("/admin")
	//管理员列表
	admin.GET("/list", handle.List)
	// 添加管理员
	admin.POST("/create", handle.Create)
	// 登录
	admin.POST("/login", handle.Login)
	// 获取当前用户信息（需登录）
	admin.GET("/profile", builtin.JWTAuthMiddleware(), handle.Profile)
}
