package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

func TestUseDefaultMiddlewaresAttachesRecoverySecurityAndRateLimit(t *testing.T) {
	resetRateLimitStore()

	engine := gin.New()
	UseDefaultMiddlewares(engine, DefaultOptions{
		RequestBodyLimit: 10,
		UploadBodyLimit:  100,
		RateLimitEnabled: true,
		RateLimitLimit:   2,
		RateLimitWindow:  time.Minute,
		LogCapture:       false,
	})

	engine.GET("/ok", func(c *gin.Context) {
		response.Success(c, gin.H{"ok": true})
	})
	engine.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})
	engine.POST("/limited", func(c *gin.Context) {
		response.Success(c, gin.H{"ok": true})
	})

	okRecorder := httptest.NewRecorder()
	okRequest, _ := http.NewRequest(http.MethodGet, "/ok", nil)
	engine.ServeHTTP(okRecorder, okRequest)
	if okRecorder.Code != http.StatusOK {
		t.Fatalf("正常请求状态码不正确: got=%d want=%d", okRecorder.Code, http.StatusOK)
	}
	if got := okRecorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("安全头未生效: got=%s want=%s", got, "nosniff")
	}

	for i := 0; i < 2; i++ {
		rateRecorder := httptest.NewRecorder()
		rateRequest, _ := http.NewRequest(http.MethodGet, "/ok", nil)
		engine.ServeHTTP(rateRecorder, rateRequest)
	}
	blockedRecorder := httptest.NewRecorder()
	blockedRequest, _ := http.NewRequest(http.MethodGet, "/ok", nil)
	engine.ServeHTTP(blockedRecorder, blockedRequest)
	if blockedRecorder.Code != http.StatusTooManyRequests {
		t.Fatalf("第三次请求应被限流: got=%d want=%d", blockedRecorder.Code, http.StatusTooManyRequests)
	}

	bodyRecorder := httptest.NewRecorder()
	bodyRequest, _ := http.NewRequest(http.MethodPost, "/limited", strings.NewReader(strings.Repeat("a", 11)))
	bodyRequest.ContentLength = 11
	engine.ServeHTTP(bodyRecorder, bodyRequest)
	if bodyRecorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("请求体限制未生效: got=%d want=%d", bodyRecorder.Code, http.StatusRequestEntityTooLarge)
	}

	panicRecorder := httptest.NewRecorder()
	panicRequest, _ := http.NewRequest(http.MethodGet, "/panic", nil)
	engine.ServeHTTP(panicRecorder, panicRequest)
	if panicRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("Recovery 未生效: got=%d want=%d", panicRecorder.Code, http.StatusInternalServerError)
	}
}
