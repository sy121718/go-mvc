package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestGetServerParsesHTTPTimeouts(t *testing.T) {
	oldV := v
	t.Cleanup(func() {
		v = oldV
	})

	cfg := viper.New()
	cfg.Set("server.port", 8080)
	cfg.Set("server.mode", "test")
	cfg.Set("server.app_name", "go-mvc")
	cfg.Set("server.read_header_timeout", "3s")
	cfg.Set("server.read_timeout", "15s")
	cfg.Set("server.write_timeout", "30s")
	cfg.Set("server.idle_timeout", "60s")
	v = cfg

	serverCfg, err := GetServer()
	if err != nil {
		t.Fatalf("获取服务配置失败: %v", err)
	}

	if serverCfg.ReadHeaderTimeout != 3*time.Second {
		t.Fatalf("ReadHeaderTimeout 不正确: got=%s want=%s", serverCfg.ReadHeaderTimeout, 3*time.Second)
	}
	if serverCfg.ReadTimeout != 15*time.Second {
		t.Fatalf("ReadTimeout 不正确: got=%s want=%s", serverCfg.ReadTimeout, 15*time.Second)
	}
	if serverCfg.WriteTimeout != 30*time.Second {
		t.Fatalf("WriteTimeout 不正确: got=%s want=%s", serverCfg.WriteTimeout, 30*time.Second)
	}
	if serverCfg.IdleTimeout != 60*time.Second {
		t.Fatalf("IdleTimeout 不正确: got=%s want=%s", serverCfg.IdleTimeout, 60*time.Second)
	}
}
