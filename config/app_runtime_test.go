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

	runtime := NewAppRuntime(serverCfg, router, httpServer)

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
