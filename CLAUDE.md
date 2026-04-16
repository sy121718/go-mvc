# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Go MVC 项目，专为快速开发设计，采用模块化架构，支持未来微服务迁移。使用 Gin 框架，当前基础设施采用显式生命周期编排。

## 常用命令

```bash
# 运行项目
go run cmd/main.go

# 构建
go build -o app cmd/main.go

# 运行测试
go test ./...

# 运行指定包的测试
go test ./internal/module/backend/...
```

## 架构设计

### 项目目录

```text
your-project/
├── cmd/main.go                            # 启动入口
├── config/                                # 配置解析与生命周期编排
├── internal/                              # 私有应用代码（Go 编译器保护）
│   ├── routers/                           # 主路由聚合（外部整合层）
│   ├── middleware/                        # 全局中间件
│   ├── module/                            # 业务模块
│   │   └── backend/                       # 当前已启用的后台业务模块
│   └── task/                              # 任务处理
├── pkg/                                   # 通用组件
├── public/                                # 公共文件聚合目录
│   ├── docs/                              # 项目文档
│   ├── logs/                              # 运行日志目录
│   ├── storage/                           # 静态资源目录
│   └── test/                              # 测试资源目录
└── config.yaml                            # 配置文件
```

### 目录职责

- **cmd/** - 启动入口（main.go）
- **config/** - 配置解析与生命周期编排入口
- **pkg/** - 通用基础组件与能力封装
- **internal/** - 私有应用代码（Go 编译器保护）
  - **routers/** - 主路由聚合（注册所有模块路由）
  - **middleware/** - 全局中间件（auth, cors, logging, rate limit）
  - **module/** - 业务模块
  - **task/** - 异步任务处理
- **public/** - 公共文件与运行产物聚合目录
  - **docs/** - 文档资料
  - **logs/** - 运行日志
  - **storage/** - 业务静态资源
  - **test/** - 测试资源

### 模块结构

每个业务模块至少保持以下核心层级：

```text
module_name/
├── router/               # 路由目录
│   └── admin_router.go   # 按模块名命名
├── handle/               # 控制器目录
│   └── admin_handle.go
├── service/              # 业务逻辑目录
│   └── admin_service.go
└── model/                # 数据模型目录
    └── admin_model.go
```

可按需要补充：`dto/`、`enums/`、`client/`、`service/helper/` 等目录，但命名规则保持一致；若暂时不用，不创建空目录。

### 命名规则

- 外部整合层目录可以使用复数，例如 `internal/routers`
- 模块内部的小目录使用单数，例如 `router`、`handle`、`service`、`model`
- 模块内文件统一使用“模块名 + 分层名”命名：
  - `admin_router.go`
  - `admin_handle.go`
  - `admin_service.go`
  - `admin_model.go`
- 模块内包名也统一使用“模块名 + 分层名”命名，避免跨模块冲突：
  - `package adminrouter`
  - `package adminhandle`
  - `package adminservice`
  - `package adminmodel`
- 扩展层沿用同一规则：
  - `admin_dto.go` -> `package admindto`
  - `admin_enum.go` / `admin_error.go` -> `package adminenums`
  - `admin_client.go` -> `package adminclient`
  - `admin_helper.go` -> `package adminhelper`
- 不保留 `sys_` 前缀
- 不使用 `router.go`、`user.go` 这类缺少模块语义的文件名

### 核心设计模式

1. **显式生命周期编排**：基础组件统一由 `config.InitComponents()` / `config.CloseComponents()` 编排
2. **模块隔离**：每个模块独立维护自己的路由、控制器、服务、模型
3. **配置分离**：默认值在 `config/config.go` 里设置，各 pkg 自己解析配置段
4. **数据库驱动 i18n**：多语言配置中心以数据库为唯一数据源
5. **基础设施门面 + provider/driver 实现**：`pkg` 根入口对外暴露统一 API，具体技术栈实现放在 `provider/` 或 `driver/` 子目录

### 配置读取流程

```text
config.yaml → config/config.go（默认值 + 解析）→ pkg 根入口（database/cache/queue 等）→ provider/driver 具体实现
```

默认值应在 `config/config.go` 中通过 `viper.SetDefault()` 设置，不在 pkg 组件中设置。

### 日志与 public 目录

- 日志文件默认放在 `public/logs/app.log`
- `public/logs/` 与 `public/storage/` 分工不同：
  - `logs/` 用于运行日志与排查产物
  - `storage/` 用于业务静态资源
- `public/` 只聚合公共文件与运行产物，不放业务实现代码

## 开发规范

- 文档和注释使用中文
- Windows 开发环境 - 使用 Windows 路径和命令
- 不生成测试脚本 - 用户自行测试
- 不自动启动项目 - 用户手动启动
- 使用 Gin binding 标签做请求验证（不单独建 validator 层）
- 涉及 internal/module 下的业务模块开发时，先阅读 internal/module/CLAUDE.md，遵循模块结构与多语言规范
- 涉及 public 下的文档、静态资源、测试资源、日志目录时，先阅读 public/CLAUDE.md，遵循目录边界

### 路由规范

**路由分层：**

- **主路由**：`internal/routers/routes.go` - 聚合所有模块路由
- **模块路由**：`internal/module/{分组}/{模块}/router/{模块}_router.go` - 定义模块内路由

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
