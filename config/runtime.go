package config

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"go-mvc/pkg/enums"
	"go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
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

	cfg := GetViper()
	prepareRuntimeRegistrations()
	components := registeredComponents()
	initialized := make([]runtimeComponent, 0, len(components))

	log.Println("开始初始化组件...")
	for _, component := range components {
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

	router.GET("/health", func(c *gin.Context) {
		response.SuccessWithMessage(c, enums.MsgOperationSuccess, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api")
	for _, module := range registeredModules() {
		module.Register(api)
	}

	router.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "请求的资源不存在")
	})
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
