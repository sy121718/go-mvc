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
		auth.GET("/list", handle.List)
		auth.GET("/detail", handle.Detail)
		auth.POST("/create", handle.Create)
		auth.GET("/profile", handle.Profile)
		auth.POST("/edit", handle.Edit)
	}
}
