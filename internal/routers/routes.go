package routers

import (
	"net/http"

	adminhttp "go-mvc/internal/module/backend/admin/inbound/http"
	captcharouter "go-mvc/internal/module/common/captcha/router"
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
	captcharouter.SetupCaptchaRoutes(api)
	adminhttp.SetupAdminRoutes(api)

	router.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "请求的资源不存在")
	})
}
