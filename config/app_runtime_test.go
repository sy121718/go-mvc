package config

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewAppRuntimeHoldsCoreRuntimeObjects(t *testing.T) {
	serverCfg := ServerConfig{Port: 8088}
	router := gin.New()
	httpServer := &http.Server{Addr: ":8088", Handler: router}

	runtime := NewAppRuntime(serverCfg, router, httpServer, nil, nil)

	if runtime.Server.Port != 8088 {
		t.Fatalf("ServerConfig 不正确: got=%d want=%d", runtime.Server.Port, 8088)
	}
	if runtime.Router != router {
		t.Fatalf("Router 未正确挂载")
	}
	if runtime.HTTPServer != httpServer {
		t.Fatalf("HTTPServer 未正确挂载")
	}
}

func TestAppRuntimeCanCarryReadyAndShutdownFunctions(t *testing.T) {
	calledReady := false
	calledShutdown := false

	runtime := NewAppRuntime(
		ServerConfig{Port: 8088},
		nil,
		nil,
		func() error {
			calledReady = true
			return nil
		},
		func() error {
			calledShutdown = true
			return nil
		},
	)

	if err := runtime.Ready(); err != nil {
		t.Fatalf("Ready 执行失败: %v", err)
	}
	if err := runtime.Shutdown(); err != nil {
		t.Fatalf("Shutdown 执行失败: %v", err)
	}
	if !calledReady || !calledShutdown {
		t.Fatalf("运行时回调未被正确挂载")
	}
}
