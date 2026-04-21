package main

import (
	"context"
	"errors"
	"fmt"
	"go-mvc/config"
	"go-mvc/internal/middleware"
	"go-mvc/internal/routers"
	"go-mvc/pkg/utils"
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
	// 1) 读取配置文件，包含 server/database/redis 等全局配置。
	if err := config.Init("config.yaml"); err != nil {
		return fmt.Errorf("配置加载失败: %w", err)
	}

	// 2) 解析服务配置并设置 Gin 运行模式（debug/release/test）。
	serverCfg, err := config.GetServer()
	if err != nil {
		return err
	}
	gin.SetMode(serverCfg.Mode)

	// 3) 初始化基础组件（DB、i18n、cache、auth、upload、queue...）。
	if err := config.InitComponents(); err != nil {
		return fmt.Errorf("组件初始化失败: %w", err)
	}

	// 4) 构建 HTTP 路由。
	router := buildHTTPRouter(
		serverCfg,
		config.GetViper().GetBool("log.capture.http"),
		config.ModuleRegistrars(),
		config.ValidateReady,
	)

	// 5) 启动前端口策略：
	// - debug/test：自动尝试释放端口（仅白名单进程）
	// - release：只提示占用进程与 kill 命令，不自动结束进程
	if err := utils.EnsurePortReady(serverCfg.Mode, serverCfg.Port); err != nil {
		return err
	}

	addr := fmt.Sprintf(":%d", serverCfg.Port)
	log.Printf("服务启动: http://localhost%s", addr)

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: serverCfg.ReadHeaderTimeout,
		ReadTimeout:       serverCfg.ReadTimeout,
		WriteTimeout:      serverCfg.WriteTimeout,
		IdleTimeout:       serverCfg.IdleTimeout,
	}

	// 6) 在 goroutine 里启动 HTTP 服务：
	// - 正常关闭时 ListenAndServe 会返回 http.ErrServerClosed（不是错误）
	// - 非预期错误（如端口冲突）通过 channel 回传到主协程统一处理
	serverErrCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	// 7) 等待两类信号：
	// - 服务启动/运行期错误
	// - 进程退出信号（Ctrl+C / kill）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case err := <-serverErrCh:
		if err != nil {
			// 启动失败时也要执行组件关闭，避免留下后台资源（连接、定时器等）。
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

	// 8) 关闭 HTTP：给在途请求最多 5 秒收尾时间。
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP Server 关闭失败: %v", err)
	}

	// 9) 统一关闭组件，释放数据库、缓存、队列等资源。
	if err := config.CloseComponents(); err != nil {
		log.Printf("组件关闭失败: %v", err)
	}

	log.Println("服务已退出")
	return nil
}

func buildHTTPRouter(
	serverCfg config.ServerConfig,
	logCapture bool,
	modules []func(*gin.RouterGroup),
	ready func() error,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.RequestBodyLimitMiddleware(serverCfg.RequestBodyLimit, serverCfg.UploadBodyLimit))
	if serverCfg.RateLimitEnabled {
		router.Use(middleware.RequestRateLimitMiddleware(serverCfg.RateLimitLimit, serverCfg.RateLimitWindow))
	}
	router.Use(middleware.RequestLogCaptureMiddleware(logCapture))
	routers.SetupRoutes(router, modules, ready)
	return router
}
