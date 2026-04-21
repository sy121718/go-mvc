package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-mvc/pkg/enums"
	"go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

func TestSetupRoutesRegistersBaseRoutesAndModules(t *testing.T) {
	engine := gin.New()
	SetupRoutes(engine, []func(*gin.RouterGroup){
		func(rg *gin.RouterGroup) {
			rg.GET("/mock", func(c *gin.Context) {
				response.Success(c, gin.H{"module": "ok"})
			})
		},
	}, func() error {
		return nil
	})

	livezRecorder := httptest.NewRecorder()
	livezRequest, _ := http.NewRequest(http.MethodGet, "/livez", nil)
	engine.ServeHTTP(livezRecorder, livezRequest)
	if livezRecorder.Code != http.StatusOK {
		t.Fatalf("/livez 状态码不正确: got=%d want=%d", livezRecorder.Code, http.StatusOK)
	}

	moduleRecorder := httptest.NewRecorder()
	moduleRequest, _ := http.NewRequest(http.MethodGet, "/api/mock", nil)
	engine.ServeHTTP(moduleRecorder, moduleRequest)
	if moduleRecorder.Code != http.StatusOK {
		t.Fatalf("/api/mock 状态码不正确: got=%d want=%d", moduleRecorder.Code, http.StatusOK)
	}

	notFoundRecorder := httptest.NewRecorder()
	notFoundRequest, _ := http.NewRequest(http.MethodGet, "/not-found", nil)
	engine.ServeHTTP(notFoundRecorder, notFoundRequest)

	var notFoundResp response.Response
	if err := json.Unmarshal(notFoundRecorder.Body.Bytes(), &notFoundResp); err != nil {
		t.Fatalf("解析 404 响应失败: %v", err)
	}
	if notFoundResp.Code != enums.ErrNotFound {
		t.Fatalf("404 错误码不正确: got=%s want=%s", notFoundResp.Code, enums.ErrNotFound)
	}
}

func TestSetupRoutesReturns503WhenReadyCheckFails(t *testing.T) {
	engine := gin.New()
	SetupRoutes(engine, nil, func() error {
		return fmt.Errorf("runtime not initialized")
	})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/readyz", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("/readyz 状态码不正确: got=%d want=%d", recorder.Code, http.StatusServiceUnavailable)
	}

	var resp response.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析 /readyz 响应失败: %v", err)
	}
	if resp.Code != "ErrNotReady" {
		t.Fatalf("/readyz 错误码不正确: got=%s want=%s", resp.Code, "ErrNotReady")
	}
}
