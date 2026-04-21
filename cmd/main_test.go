package main

import (
	"net/http"
	"testing"
	"time"

	"go-mvc/config"
	"go-mvc/public/test/support"

	"github.com/gin-gonic/gin"
)

func TestBuildHTTPRouterUsesRecoveryAndDefaultMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := buildHTTPRouter(
		config.ServerConfig{
			RequestBodyLimit: 10,
			UploadBodyLimit:  100,
			RateLimitEnabled: true,
			RateLimitLimit:   2,
			RateLimitWindow:  time.Minute,
		},
		false,
		[]func(*gin.RouterGroup){
			func(rg *gin.RouterGroup) {
				rg.GET("/panic", func(c *gin.Context) {
					panic("boom")
				})
			},
		},
		func() error { return nil },
	)

	recorder, err := support.SendRequest(router, support.RequestOptions{
		Method: http.MethodGet,
		Path:   "/livez",
	})
	if err != nil {
		t.Fatalf("请求 /livez 失败: %v", err)
	}
	if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("安全头未生效: got=%s want=%s", got, "nosniff")
	}

	panicRecorder, err := support.SendRequest(router, support.RequestOptions{
		Method: http.MethodGet,
		Path:   "/api/panic",
	})
	if err != nil {
		t.Fatalf("请求 panic 路由失败: %v", err)
	}
	if panicRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("Recovery 未生效: got=%d want=%d", panicRecorder.Code, http.StatusInternalServerError)
	}
}

func TestBuildHTTPServerAppliesServerConfig(t *testing.T) {
	handler := http.NewServeMux()
	serverCfg := config.ServerConfig{
		Port:              8088,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	server := buildHTTPServer(serverCfg, handler)

	if server.Addr != ":8088" {
		t.Fatalf("Addr 不正确: got=%s want=%s", server.Addr, ":8088")
	}
	if server.Handler != handler {
		t.Fatalf("Handler 未正确挂载")
	}
	if server.ReadHeaderTimeout != 3*time.Second {
		t.Fatalf("ReadHeaderTimeout 不正确: got=%s want=%s", server.ReadHeaderTimeout, 3*time.Second)
	}
	if server.ReadTimeout != 15*time.Second {
		t.Fatalf("ReadTimeout 不正确: got=%s want=%s", server.ReadTimeout, 15*time.Second)
	}
	if server.WriteTimeout != 30*time.Second {
		t.Fatalf("WriteTimeout 不正确: got=%s want=%s", server.WriteTimeout, 30*time.Second)
	}
	if server.IdleTimeout != 60*time.Second {
		t.Fatalf("IdleTimeout 不正确: got=%s want=%s", server.IdleTimeout, 60*time.Second)
	}
}
