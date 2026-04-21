package routers

import (
	"net/http"

	"go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 装配系统基础路由和业务模块路由。
//
// 参数说明：
// - router: Gin 引擎实例
// - modules: 业务模块路由注册函数列表，统一挂到 /api 分组下
// - ready: 就绪检查函数，用于 /readyz 返回服务是否可提供服务
func SetupRoutes(
	router *gin.Engine,
	modules []func(*gin.RouterGroup),
	ready func() error,
) {
	if router == nil {
		return
	}

	router.GET("/livez", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Response{
			Code:    "0",
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
					Code:    "ErrNotReady",
					Message: err.Error(),
					Data: gin.H{
						"status": "not_ready",
					},
				})
				return
			}
		}

		c.JSON(http.StatusOK, response.Response{
			Code:    "0",
			Message: "ok",
			Data: gin.H{
				"status": "ready",
			},
		})
	})

	api := router.Group("/api")
	for _, register := range modules {
		register(api)
	}

	router.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "请求的资源不存在")
	})
}
