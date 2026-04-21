package config

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

var (
	runtimeMu           sync.Mutex
	initializedRegistry []runtimeComponent
	runtimeInited       bool
)

const (
	errComponentInitFailedPrefix = "component init failed"
	errComponentNotReadyPrefix   = "component not ready"
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
	for _, critical := range []bool{true, false} {
		for _, component := range runtimeComponents {
			if component.Critical != critical {
				continue
			}
			if component.Enabled != nil && !component.Enabled(cfg) {
				continue
			}

			if err := component.Init(cfg); err != nil {
				_ = closeComponents(initialized)
				return fmt.Errorf("%s [%s]: %w", errComponentInitFailedPrefix, component.Name, err)
			}
			initialized = append(initialized, component)
		}
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

// ValidateReady 检查当前运行时是否达到“可对外提供服务”的就绪状态。
func ValidateReady() error {
	runtimeMu.Lock()
	ready := runtimeInited
	runtimeMu.Unlock()

	if !ready {
		return fmt.Errorf("%s [runtime]: runtime not initialized", errComponentNotReadyPrefix)
	}

	for _, component := range initializedRegistry {
		if component.Ready == nil {
			continue
		}
		if err := component.Ready(); err != nil {
			return fmt.Errorf("%s [%s]: %w", errComponentNotReadyPrefix, component.Name, err)
		}
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
