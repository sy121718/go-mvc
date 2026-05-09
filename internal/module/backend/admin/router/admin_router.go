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

	admin := rg.Group("/admin")
	admin.GET("/list", handle.List)
	admin.POST("/create", handle.Create)
}
