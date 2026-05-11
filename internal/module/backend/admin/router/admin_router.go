package adminrouter

import "github.com/gin-gonic/gin"

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
}
