package routers

import (
	"net/http"

	adminrouter "go-mvc/internal/module/backend/admin/router"
	userrouter "go-mvc/internal/module/backend/user/router"
	"go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, ready func() error) {
	if router == nil {
		return
	}

	router.GET("/livez", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Response{
			Code:    http.StatusOK,
			Message: "ok",
			Data: gin.H{
				"status": "alive",
			},
		})
	})

	router.GET("/readyz", func(c *gin.Context) {
		if ready != nil {
			if err := ready(); err != nil {
				c.JSON(http.StatusServiceUnavailable, response.Response{
					Code:    http.StatusServiceUnavailable,
					Message: err.Error(),
					Data: gin.H{
						"status": "not_ready",
					},
				})
				return
			}
		}

		c.JSON(http.StatusOK, response.Response{
			Code:    http.StatusOK,
			Message: "ok",
			Data: gin.H{
				"status": "ready",
			},
		})
	})

	api := router.Group("/api")
	adminrouter.SetupAdminRoutes(api)
	userrouter.SetupUserRoutes(api)

	router.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "请求的资源不存在")
	})
}
