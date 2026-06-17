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
	admin.POST("/login", handle.Login)

	// 以下路由需要 JWT 鉴权
	auth := admin.Group("").Use(builtin.JWTAuthMiddleware())
	{
		//列表
		auth.GET("/list", handle.List)
		//查询详情
		auth.GET("/detail", handle.Detail)
		//创建
		auth.POST("/create", handle.Create)
		// 编辑管理员信息
		auth.POST("/edit", handle.Edit)
		// 个人信息
		auth.GET("/profile", handle.Profile)
		//删除管理员
	}
}
