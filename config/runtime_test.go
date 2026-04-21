package config

import (
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestInitComponentsAndCloseComponentsFollowRegisteredOrder(t *testing.T) {
	oldV := v
	oldRuntimeInited := runtimeInited
	oldInitializedRegistry := initializedRegistry
	oldPreparers := runtimePreparers
	oldComponents := runtimeComponents

	t.Cleanup(func() {
		v = oldV
		runtimeInited = oldRuntimeInited
		initializedRegistry = oldInitializedRegistry
		runtimePreparers = oldPreparers
		runtimeComponents = oldComponents
	})

	v = viper.New()
	runtimeInited = false
	initializedRegistry = nil

	events := make([]string, 0, 8)
	runtimePreparers = []func(){
		func() {
			events = append(events, "prepare")
		},
	}
	runtimeComponents = []runtimeComponent{
		{
			Name:     "first",
			Critical: true,
			Init: func(_ *viper.Viper) error {
				events = append(events, "init:first")
				return nil
			},
			Close: func() error {
				events = append(events, "close:first")
				return nil
			},
		},
		{
			Name: "disabled",
			Enabled: func(_ *viper.Viper) bool {
				return false
			},
			Init: func(_ *viper.Viper) error {
				events = append(events, "init:disabled")
				return nil
			},
			Close: func() error {
				events = append(events, "close:disabled")
				return nil
			},
		},
		{
			Name: "second",
			Init: func(_ *viper.Viper) error {
				events = append(events, "init:second")
				return nil
			},
			Close: func() error {
				events = append(events, "close:second")
				return nil
			},
		},
	}

	if err := InitComponents(); err != nil {
		t.Fatalf("初始化组件失败: %v", err)
	}

	if err := CloseComponents(); err != nil {
		t.Fatalf("关闭组件失败: %v", err)
	}

	want := []string{
		"prepare",
		"init:first",
		"init:second",
		"close:second",
		"close:first",
	}
	if !reflect.DeepEqual(events, want) {
		t.Fatalf("组件编排顺序不正确:\nwant=%v\ngot=%v", want, events)
	}
}

func TestInitComponentsInitializesCriticalBeforeOptional(t *testing.T) {
	oldV := v
	oldRuntimeInited := runtimeInited
	oldInitializedRegistry := initializedRegistry
	oldPreparers := runtimePreparers
	oldComponents := runtimeComponents

	t.Cleanup(func() {
		v = oldV
		runtimeInited = oldRuntimeInited
		initializedRegistry = oldInitializedRegistry
		runtimePreparers = oldPreparers
		runtimeComponents = oldComponents
	})

	v = viper.New()
	v.Set("server.mode", "test")
	runtimeInited = false
	initializedRegistry = nil
	runtimePreparers = nil

	events := make([]string, 0, 4)
	runtimeComponents = []runtimeComponent{
		{
			Name: "optional",
			Init: func(_ *viper.Viper) error {
				events = append(events, "init:optional")
				return nil
			},
		},
		{
			Name:     "critical",
			Critical: true,
			Init: func(_ *viper.Viper) error {
				events = append(events, "init:critical")
				return nil
			},
		},
	}

	if err := InitComponents(); err != nil {
		t.Fatalf("初始化组件失败: %v", err)
	}

	want := []string{"init:critical", "init:optional"}
	if !reflect.DeepEqual(events, want) {
		t.Fatalf("关键组件应先初始化:\nwant=%v\ngot=%v", want, events)
	}
}

func TestInitComponentsReturnsClassifiedInitError(t *testing.T) {
	oldV := v
	oldRuntimeInited := runtimeInited
	oldInitializedRegistry := initializedRegistry
	oldPreparers := runtimePreparers
	oldComponents := runtimeComponents

	t.Cleanup(func() {
		v = oldV
		runtimeInited = oldRuntimeInited
		initializedRegistry = oldInitializedRegistry
		runtimePreparers = oldPreparers
		runtimeComponents = oldComponents
	})

	v = viper.New()
	v.Set("server.mode", "test")
	runtimeInited = false
	initializedRegistry = nil
	runtimePreparers = nil
	runtimeComponents = []runtimeComponent{
		{
			Name: "cache",
			Init: func(_ *viper.Viper) error {
				return errString("dial tcp failed")
			},
		},
	}

	err := InitComponents()
	if err == nil {
		t.Fatalf("组件初始化失败时应返回错误")
	}
	if !strings.Contains(err.Error(), "component init failed") {
		t.Fatalf("错误分类不正确: got=%v", err)
	}
	if !strings.Contains(err.Error(), "cache") {
		t.Fatalf("错误中应包含组件名: got=%v", err)
	}
}

func TestValidateReadyReturnsClassifiedReadyError(t *testing.T) {
	oldV := v
	oldRuntimeInited := runtimeInited
	oldInitializedRegistry := initializedRegistry

	t.Cleanup(func() {
		v = oldV
		runtimeInited = oldRuntimeInited
		initializedRegistry = oldInitializedRegistry
	})

	v = viper.New()
	v.Set("server.mode", "test")
	runtimeInited = true
	initializedRegistry = []runtimeComponent{
		{
			Name: "queue",
			Ready: func() error {
				return errString("not ready")
			},
		},
	}

	err := ValidateReady()
	if err == nil {
		t.Fatalf("组件未就绪时应返回错误")
	}
	if !strings.Contains(err.Error(), "component not ready") {
		t.Fatalf("错误分类不正确: got=%v", err)
	}
	if !strings.Contains(err.Error(), "queue") {
		t.Fatalf("错误中应包含组件名: got=%v", err)
	}
}

type errString string

func (e errString) Error() string {
	return string(e)
}
