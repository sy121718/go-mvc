package adminhttp

import (
	"go-mvc/internal/middleware/builtin"
	adminmodel "go-mvc/internal/module/backend/admin/model"
	adminservice "go-mvc/internal/module/backend/admin/service"
	"go-mvc/pkg/database"

	"github.com/gin-gonic/gin"
)

func newHandle() (*Handle, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}
	service := adminservice.NewService(adminmodel.NewAdminModel(db))
	return NewHandle(service), nil
}

func SetupAdminRoutes(rg *gin.RouterGroup) {
	if rg == nil {
		return
	}

	handle, err := newHandle()
	if err != nil {
		return
	}

	admin := rg.Group("/admin")
	admin.POST("/login", handle.Login)

	auth := admin.Group("").Use(builtin.JWTAuthMiddleware())
	{
		auth.GET("/list", handle.List)
		auth.GET("/detail", handle.Detail)
		auth.POST("/create", handle.Create)
		auth.POST("/edit", handle.Edit)
		auth.GET("/profile", handle.Profile)
	}
}
