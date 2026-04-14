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
	// 1. 加载配置文件
	if err := config.Init("config.yaml"); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 获取 viper 实例
	v := config.GetViper()

	// 设置 Gin 模式
	serverCfg := config.GetServer()
	gin.SetMode(serverCfg.Mode)

	// 2. 初始化组件（config 驱动 pkg 初始化，自动监听退出信号）
	config.InitComponents(v)

	// 3. 创建 Gin 引擎
	router := gin.Default()

	// 4. 注册路由
	routers.SetupRoutes(router)

	// 5. 启动服务
	addr := fmt.Sprintf(":%d", serverCfg.Port)
	log.Printf("服务启动: http://localhost%s", addr)

	// 在 goroutine 中启动服务器
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("收到退出信号，开始关闭...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. 先关闭 HTTP Server（停止接收新请求）
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP Server 关闭失败: %v", err)
	}

	// 2. 再关闭数据库和 Redis
	config.CloseComponents()

	log.Println("服务已退出")
}
