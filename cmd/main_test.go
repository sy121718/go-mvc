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
