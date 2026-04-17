package main

import (
	"context"
	"errors"
	"fmt"
	"go-mvc/config"
	"go-mvc/internal/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

/*
Go 程序入口 - main 函数
===========================================
PHP 对比：
- PHP: 入口文件通常是 index.php，通过 Web 服务器访问
- Go: main 函数是程序唯一入口，编译后直接运行

启动流程：
1. 加载配置文件
2. 初始化组件（config 驱动 pkg 初始化）
3. 创建 Gin 引擎
4. 注册路由
5. 启动 HTTP 服务

配置管理说明：
- config.Init() 读取配置文件，返回 viper 实例
- config.InitComponents() 驱动所有 pkg 组件初始化
- 各个 pkg 定义自己的 Config 结构体和默认值
*/
func main() {
	if err := run(); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

func run() error {
	if err := config.Init("config.yaml"); err != nil {
		return fmt.Errorf("配置加载失败: %w", err)
	}

	serverCfg, err := config.GetServer()
	if err != nil {
		return err
	}
	gin.SetMode(serverCfg.Mode)

	if err := config.InitComponents(); err != nil {
		return fmt.Errorf("组件初始化失败: %w", err)
	}

	router := gin.Default()
	routers.SetupRoutes(router)

	addr := fmt.Sprintf(":%d", serverCfg.Port)
	log.Printf("服务启动: http://localhost%s", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case err := <-serverErrCh:
		if err != nil {
			if closeErr := config.CloseComponents(); closeErr != nil {
				log.Printf("组件关闭失败: %v", closeErr)
			}
			return fmt.Errorf("HTTP 服务启动失败: %w", err)
		}
		return nil
	case <-quit:
		log.Println("收到退出信号，开始关闭...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP Server 关闭失败: %v", err)
	}

	if err := config.CloseComponents(); err != nil {
		log.Printf("组件关闭失败: %v", err)
	}

	log.Println("服务已退出")
	return nil
}
