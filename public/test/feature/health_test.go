package feature

import (
	"net/http"
	"testing"

	"go-mvc/public/test/support"
)

func TestHealthCheckSuccess(t *testing.T) {
	engine, cleanup, err := support.SetupTestBootstrap(support.BootstrapOptions{
		UseDefaultRoute: true,
		InitComponents:  false,
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
		Method: http.MethodGet,
		Path:   "/health",
	})
	if err != nil {
		t.Fatalf("发送请求失败: %v", err)
	}

	if recorder.Code != http.StatusOK {
		t.Fatalf("状态码不正确: got=%d want=%d", recorder.Code, http.StatusOK)
	}

	response, err := support.ParseStandardResponse(recorder)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response.Code != "0" {
		t.Fatalf("响应 code 不正确: got=%s want=%s", response.Code, "0")
	}

	var data struct {
		Status string `json:"status"`
	}
	if err := support.DecodeResponseData(recorder, &data); err != nil {
		t.Fatalf("解析 data 失败: %v", err)
	}

	if data.Status != "ok" {
		t.Fatalf("响应 data.status 不正确: got=%s want=%s", data.Status, "ok")
	}
}
