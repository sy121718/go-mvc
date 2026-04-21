package config

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"go-mvc/pkg/auth"
	"go-mvc/pkg/cache"
	"go-mvc/pkg/casbin"
	"go-mvc/pkg/database"
	"go-mvc/pkg/i18n"
	"go-mvc/pkg/queue"
	"go-mvc/pkg/response"
	"go-mvc/pkg/upload"

	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	runtimeMu           sync.Mutex
	initializedRegistry []runtimeComponent
	runtimeInited       bool
)

// InitComponents 按注册顺序初始化运行时组件。
func InitComponents() error {
	runtimeMu.Lock()
	defer runtimeMu.Unlock()

	if runtimeInited {
		return nil
	}

	if err := ValidateRuntimeConfig(); err != nil {
		return err
	}

	cfg := GetViper()
	for _, prepare := range runtimePreparers {
		prepare()
	}

	initialized := make([]runtimeComponent, 0, len(runtimeComponents))

	log.Println("开始初始化组件...")
	for _, component := range runtimeComponents {
		if component.Enabled != nil && !component.Enabled(cfg) {
			continue
		}

		if err := component.Init(cfg); err != nil {
			_ = closeComponents(initialized)
			return fmt.Errorf("初始化组件 %s 失败: %w", component.Name, err)
		}
		initialized = append(initialized, component)
	}

	initializedRegistry = initialized
	runtimeInited = true
	log.Println("组件初始化完成")
	return nil
}

// CloseComponents 按初始化逆序关闭运行时组件。
func CloseComponents() error {
	runtimeMu.Lock()
	defer runtimeMu.Unlock()

	if !runtimeInited {
		return nil
	}

	log.Println("开始关闭组件...")
	closeErr := closeComponents(initializedRegistry)
	initializedRegistry = nil
	runtimeInited = false

	if closeErr != nil {
		return closeErr
	}

	log.Println("组件关闭完成")
	return nil
}

// SetupRoutes 根据注册清单装配系统路由与业务模块路由。
func SetupRoutes(router *gin.Engine) {
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
		if err := ValidateReady(); err != nil {
			c.JSON(http.StatusServiceUnavailable, response.Response{
				Code:    "ErrNotReady",
				Message: err.Error(),
				Data: gin.H{
					"status": "not_ready",
				},
			})
			return
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
	for _, module := range runtimeModules {
		module.Register(api)
	}

	router.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "请求的资源不存在")
	})
}

// ValidateReady 检查当前运行时是否达到“可对外提供服务”的就绪状态。
func ValidateReady() error {
	runtimeMu.Lock()
	ready := runtimeInited
	runtimeMu.Unlock()

	if !ready {
		return fmt.Errorf("runtime not initialized")
	}

	cfg := GetViper()
	if !database.IsInited() {
		return fmt.Errorf("database not ready")
	}
	if err := auth.MustBeReady(); err != nil {
		return fmt.Errorf("auth not ready: %w", err)
	}
	if err := i18n.ValidateReady(); err != nil {
		return fmt.Errorf("i18n not ready: %w", err)
	}
	if cfg.GetBool("redis.enabled") && !cache.IsInited() {
		return fmt.Errorf("cache not ready")
	}
	if cfg.GetBool("casbin.enabled") && casbin.GetEnforcer() == nil {
		return fmt.Errorf("casbin not ready")
	}
	if cfg.GetBool("upload.enabled") && !upload.IsInited() {
		return fmt.Errorf("upload not ready")
	}
	if cfg.GetBool("queue.enabled") && !queue.IsInited() {
		return fmt.Errorf("queue not ready")
	}
	return nil
}

func closeComponents(components []runtimeComponent) error {
	var closeErr error
	for i := len(components) - 1; i >= 0; i-- {
		component := components[i]
		if component.Close == nil {
			continue
		}
		if err := component.Close(); err != nil {
			closeErr = errors.Join(closeErr, fmt.Errorf("关闭组件 %s 失败: %w", component.Name, err))
		}
	}
	return closeErr
}
