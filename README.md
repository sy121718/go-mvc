# Go MVC

快速开发的 Go MVC 项目，采用模块化架构，支持未来微服务迁移。

## 特性

- 基于 Gin 框架
- 模块化架构设计
- 显式生命周期初始化
- 支持 MySQL、Redis、JWT、Casbin、i18n
- 任务队列支持（Asynq）

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/sy121718/go-mvc.git
cd go-mvc
```

### 2. 配置

```bash
cp config.yaml.example config.yaml
# 编辑 config.yaml，填入你的配置
```

其中 `i18n.default_lang` 用于配置网站默认语言；当请求头和 `lang` 参数都未传时，统一回退到这个语言。

### 3. 安装依赖

```bash
go mod download
```

### 4. 运行

```bash
go run cmd/main.go
```

## 项目结构

```text
go-mvc/
├── cmd/main.go                            # 启动入口
├── config/                                # 配置解析与生命周期编排
├── internal/                              # 私有应用代码
│   ├── routers/                           # 主路由聚合（外部整合层）
│   ├── middleware/                        # 全局中间件
│   ├── module/                            # 业务模块
│   │   └── backend/
│   │       ├── admin/
│   │       │   ├── router/                # 模块路由目录
│   │       │   │   └── admin_router.go
│   │       │   ├── handle/                # 控制器目录
│   │       │   │   └── admin_handle.go
│   │       │   ├── service/               # 业务层目录
│   │       │   │   └── admin_service.go
│   │       │   └── model/                 # 模型层目录
│   │       │       └── admin_model.go
│   │       └── user/
│   │           ├── router/
│   │           │   └── user_router.go
│   │           ├── handle/
│   │           │   └── user_handle.go
│   │           ├── service/
│   │           │   └── user_service.go
│   │           └── model/
│   │               └── user_model.go
│   └── task/                              # 任务处理
├── pkg/                                   # 通用组件
├── public/                                # 公共文件聚合目录
│   ├── docs/                              # 项目文档
│   ├── logs/                              # 运行日志目录
│   ├── storage/                           # 静态文件目录
│   └── test/                              # 单元测试与用例测试目录
└── config.yaml                            # 配置文件
```

## 启动与生命周期

当前项目采用显式生命周期编排：

```text
config.yaml → config.Init() → config.InitComponents() → 路由注册 → HTTP 服务启动
```

其中 `i18n.default_lang` 会在 `config.InitComponents()` 阶段注入到 i18n 组件，作为网站默认语言；pkg 内部保留 `zh-CN` 兜底。

关闭顺序：

```text
HTTP Server → i18n 刷新器 → queue → Redis → DB
```

## 命名规范

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
- 不保留 `sys_` 前缀，不使用 `router.go`、`user.go` 这类缺少模块语义的文件名，可以自己用简短命名不一定参考数据库表名

## public 目录约定

- `public/` 用于聚合公共文件与运行产物
- `public/docs/` 存放项目文档
- `public/logs/` 存放运行日志
- `public/storage/` 存放业务静态资源
- `public/test/` 存放测试资源
- `public/` 下不放业务实现代码

## 开发规范

详见 [CLAUDE.md](./CLAUDE.md)

## License

MIT
