package feature

import (
	"net/http"
	"strings"
	"testing"

	"go-mvc/internal/middleware"
	"go-mvc/pkg/response"
	"go-mvc/public/test/support"

	"github.com/gin-gonic/gin"
)

func TestRequestBodyLimitRejectsLargeJSONBody(t *testing.T) {
	engine, cleanup, err := support.SetupTestBootstrap(support.BootstrapOptions{
		UseDefaultRoute: false,
		InitComponents:  false,
		RouteRegistrar: func(engine *gin.Engine) {
			engine.POST("/limited", middleware.RequestBodyLimitMiddleware(10, 100), func(c *gin.Context) {
				response.Success(c, gin.H{"reached": true})
			})
		},
	})
	if err != nil {
		t.Fatalf("初始化测试引擎失败: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := cleanup(); closeErr != nil {
			t.Errorf("清理测试资源失败: %v", closeErr)
		}
	})

	recorder, err := support.SendRequest(engine, support.RequestOptions{
		Method:  http.MethodPost,
		Path:    "/limited",
		RawBody: []byte(strings.Repeat("a", 11)),
	})
	if err != nil {
		t.Fatalf("发送请求失败: %v", err)
	}

	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("状态码不正确: got=%d want=%d", recorder.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestRequestBodyLimitUsesUploadLimitForUploadRoute(t *testing.T) {
	engine, cleanup, err := support.SetupTestBootstrap(support.BootstrapOptions{
		UseDefaultRoute: false,
		InitComponents:  false,
		RouteRegistrar: func(engine *gin.Engine) {
			engine.POST("/upload/file", middleware.RequestBodyLimitMiddleware(10, 100), func(c *gin.Context) {
				response.Success(c, gin.H{"reached": true})
			})
		},
	})
	if err != nil {
		t.Fatalf("初始化测试引擎失败: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := cleanup(); closeErr != nil {
			t.Errorf("清理测试资源失败: %v", closeErr)
		}
	})

	recorder, err := support.SendRequest(engine, support.RequestOptions{
		Method:  http.MethodPost,
		Path:    "/upload/file",
		RawBody: []byte(strings.Repeat("a", 20)),
	})
	if err != nil {
		t.Fatalf("发送请求失败: %v", err)
	}

	if recorder.Code != http.StatusOK {
		t.Fatalf("状态码不正确: got=%d want=%d", recorder.Code, http.StatusOK)
	}
}
