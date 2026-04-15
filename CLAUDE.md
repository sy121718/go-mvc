# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Go MVC 项目，专为快速开发设计，采用模块化架构，支持未来微服务迁移。使用 Gin 框架，实现懒加载初始化模式。

## 常用命令

```bash
# 运行项目
go run cmd/main.go

# 构建
go build -o app cmd/main.go

# 运行测试
go test ./...

# 运行指定包的测试
go test ./internal/module/user/...
```

## 架构设计

### 项目目录

```text
your-project/
├── cmd/main.go             # 启动入口
├── config/                 # 配置解析（仅读取 config.yaml）
├── internal/               # 私有应用代码（Go 编译器保护）
│   ├── router/             # 主路由聚合
│   ├── middleware/         # 全局中间件
│   ├── module/             # 业务模块（user, order, admin, common）
│   └── client/             # 全局客户端管理器
├── pkg/                    # 通用组件（懒加载初始化）
│   ├── database/           # MySQL
│   ├── cache/              # Redis
│   ├── auth/               # JWT
│   └── crypto/             # 签名
├── storage/                # 文件存储
├── docs/                   # 项目文档
└── config.yaml             # 配置文件
```

### 目录职责

- **cmd/** - 启动入口（main.go）
- **config/** - 仅负责配置解析（通过 Viper 读取 config.yaml）
- **pkg/** - 组件初始化（支持懒加载：database, cache, auth, crypto）
- **internal/** - 私有应用代码（Go 编译器保护）
  - **router/** - 主路由聚合（注册所有模块路由）
  - **middleware/** - 全局中间件（auth, cors, logging, rate limit）
  - **module/** - 业务模块（user, order, admin, common）
  - **client/** - 全局客户端管理器（模块间调用）
- **docs/** - 项目文档（中文）

### 模块结构

每个模块（user, order, admin）遵循以下结构：

```text
module/
├── router.go       # 路由定义
├── handle/         # 控制器
├── model/          # 数据模型（模块私有）
├── service/        # 业务逻辑
│   └── helper/     # 原子化工具（formatter, validator, analyzer, matcher）
├── dto/            # 请求/响应结构
└── client/         # 对外接口定义（微服务迁移用）
```

### 核心设计模式

1. **懒加载初始化**：pkg/ 组件使用 `sync.Once`，由 config.yaml 配置控制
2. **模块隔离**：每个模块完全独立，通过 client 层实现跨模块调用
3. **配置分离**：config/ 只解析 YAML，pkg/ 负责初始化（默认值在 config.go 设置）

### 配置读取流程

```text
config.yaml → config/config.go（解析）→ pkg/database/mysql.go（初始化）
```

默认值应在 config.go 中通过 `viper.SetDefault()` 设置，不在 pkg/ 组件中设置。

## 开发规范

- 文档和注释使用中文
- Windows 开发环境 - 使用 Windows 路径和命令
- 不生成测试脚本 - 用户自行测试
- 不自动启动项目 - 用户手动启动
- 使用 Gin binding 标签做请求验证（不单独建 validator 层）
- 涉及 internal/module 下的业务模块开发时，先阅读 internal/module/CLAUDE.md，遵循模块结构与多语言规范

### 路由规范

**路由分层：**

- **主路由**：`internal/router/router.go` - 聚合所有模块路由
- **模块路由**：`internal/module/{模块}/router.go` - 定义模块内路由

**只使用 GET 和 POST 两种请求方法：**

- **GET** - 查询操作，无数据变化的请求
  - 列表查询
  - 详情查询
  - 数据导出
  - 搜索功能

- **POST** - 数据变更操作，有数据提交的请求
  - 新增数据
  - 修改数据
  - 删除数据
  - 状态变更
  - 批量操作

**示例：**

```go
// GET 请求 - 查询类
GET  /api/user/list       // 用户列表
GET  /api/user/detail     // 用户详情
GET  /api/user/export     // 导出用户

// POST 请求 - 变更类
POST /api/user/create     // 创建用户
POST /api/user/update     // 更新用户
POST /api/user/delete     // 删除用户
POST /api/user/status     // 修改状态
```

**禁止使用 PUT、DELETE、PATCH 等其他 HTTP 方法。**

## 微服务迁移

架构支持渐进式迁移：

- 每个模块可独立成为服务
- client/ 层定义对外接口
- internal/client/ 管理跨服务通信
- 详见 docs/微服务迁移指南.md
