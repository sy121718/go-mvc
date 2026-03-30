package main

import (
	"fmt"
	"go-mvc/config"
	"go-mvc/internal/routers"
	"go-mvc/pkg/auth"
	"go-mvc/pkg/cache"
	"go-mvc/pkg/casbin"
	"go-mvc/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/*
Go 程序入口 - main 函数
===========================================
PHP 对比：
- PHP: 入口文件通常是 index.php，通过 Web 服务器访问
- Go: main 函数是程序唯一入口，编译后直接运行

启动流程：
1. 加载配置文件
2. 初始化组件（根据配置决定是否立即初始化）
3. 创建 Gin 引擎
4. 注册路由
5. 启动 HTTP 服务

配置管理说明：
- config.Init() 读取配置文件，返回 viper 实例
- 各个 pkg 定义自己的 Config 结构体
- main.go 传递 viper 给 pkg，pkg 自己解析配置
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

	// 2. 初始化组件（传递 viper 实例）
	initComponents(v)

	// 3. 创建 Gin 引擎
	// Gin 有两种模式：
	// (1). gin.Default() - 自动添加 Logger(日志) + Recovery(panic恢复) 中间件
	// (2). gin.New() - 无默认中间件，需要手动按需添加
	router := gin.Default()

	// 4. 注册路由
	routers.SetupRoutes(router)

	// 5. 启动服务
	addr := fmt.Sprintf(":%d", serverCfg.Port)
	log.Printf("服务启动: http://localhost%s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// initComponents 初始化组件
func initComponents(v *viper.Viper) {
	// 检查数据库配置是否懒加载
	var dbCfg database.Config
	if err := v.UnmarshalKey("database", &dbCfg); err != nil {
		log.Fatalf("解析数据库配置失败: %v", err)
	}
	if !dbCfg.LazyInit {
		log.Println("初始化数据库...")
		if err := database.InitDB(v); err != nil {
			log.Fatalf("数据库初始化失败: %v", err)
		}

		// 数据库初始化后，初始化 Casbin（依赖 DB 连接）
		log.Println("初始化 Casbin...")
		if err := casbin.InitCasbin(database.GetDB(v)); err != nil {
			log.Fatalf("Casbin 初始化失败: %v", err)
		}
	}

	// 检查 Redis 配置是否懒加载
	var redisCfg cache.Config
	if err := v.UnmarshalKey("redis", &redisCfg); err != nil {
		log.Fatalf("解析 Redis 配置失败: %v", err)
	}
	if !redisCfg.LazyInit {
		log.Println("初始化 Redis...")
		if err := cache.InitRedis(v); err != nil {
			log.Fatalf("Redis 初始化失败: %v", err)
		}
	}

	// 检查 JWT 配置是否懒加载
	var jwtCfg auth.Config
	if err := v.UnmarshalKey("jwt", &jwtCfg); err != nil {
		log.Fatalf("解析 JWT 配置失败: %v", err)
	}
	if !jwtCfg.LazyInit {
		log.Println("初始化 JWT...")
		if err := auth.InitJWT(v); err != nil {
			log.Fatalf("JWT 初始化失败: %v", err)
		}
	}
}
