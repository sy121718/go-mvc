# CLAUDE.md

本文件描述当前仓库的实际开发约定。

## 项目概览

Go + Gin 后端 + Vue 3 / vue-pure-admin 前端，全栈 Web 项目。

目录结构：`cmd` / `config` / `internal` / `pkg` / `web` / `public`

## 常用命令

```bash
# 后端
go run cmd/main.go
go build -o app cmd/main.go
go test ./...

# 前端
cd web && pnpm dev
cd web && pnpm build
```

## 核心约定

### 启动与关闭

统一入口：

- `config.Init()` → 读配置
- `config.InitComponents()` → 初始化所有 `pkg` 组件
- `config.CloseComponents()` → 逆序关闭

组件不自行决定进程退出，组件只返回 `error`。配置校验在各自 `pkg.Init()` 内部完成。

### 认证与鉴权

三层分工：

- `JWT` → 认证（验签 + 解析 user_id，写入 gin.Context）
- `Casbin` → 鉴权（内存 enforcer，Enforce(user_id, path, method)）
- `Redis` → 运行时状态（封禁标记、在线心跳、用户信息缓存）

菜单即权限的可视化：`sys_menus` type=2,3 必含 permission_code，给角色/用户分配后写入 `sys_casbin_rule`。

已实现中间件：JWTAuthMiddleware、CasbinMiddleware、安全响应头、BodyLimit、RateLimit、Recovery。

### 路由

- 只用 `GET` 和 `POST`
- 禁止 RESTful 路径参数，全部用 Query 参数
- 主路由聚合在 `internal/routers/routes.go`
- 模块路由在 `internal/module/<模块>/inbound/http/`
- 健康检查：`GET /livez`、`GET /readyz`

### 响应与错误处理

统一响应格式：

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

- `pkg` 和系统包直接用中文错误提示或原始 `err`，不用 i18n
- 业务模块统一通过模块 `enums` 提供响应消息和错误消息
- 未接好 `i18n` 时，模块 `enums` 的值可直接写中文常量；接好后再切到 i18n key 或 i18n 取值
- 业务模块响应调用 `response.Success/ErrorWithMessage` 时，消息参数统一取自模块 `enums`
- 底层系统错误优先返回原始 `err`，不过度包装

### 测试

- 默认跑现有测试，不新增额外测试框架
- 接口优先维护 feature 链路测试，复杂逻辑再补 unit

### 模块开发

模块内部按 `contract / inbound / outbound / service / model / dto` 分层。

其中 `contract/` 统一存放两类契约：

- 本模块对外暴露契约：别人怎么调我
- 本模块对外依赖契约：我业务层依赖外部什么能力

统一约定：`service` 只依赖本模块 `contract/`；`inbound` 负责承接外部调用；`outbound` 负责实现本模块对外依赖契约。详见：

- [internal/module/CLAUDE.md](./internal/module/CLAUDE.md)

`outbound` 只负责外部调用实现（RPC / HTTP / MQ / SDK / cache），不承担业务决策。

`pkg` 组件方向详见：

- [pkg/CLAUDE.md](./pkg/CLAUDE.md)

### git 提交

每次 commit 用中文注明改动文件路径和修改内容简述。

## 目录职责

| 目录 | 职责 |
|------|------|
| `cmd/` | 启动入口 |
| `config/` | 配置读取、组件注册与关闭编排 |
| `internal/middleware/` | 全局中间件 |
| `internal/routers/` | 主路由聚合 |
| `internal/module/` | 业务模块 |
| `internal/task/` | 任务注册与调度 |
| `pkg/` | 可复用基础组件（facade + provider/driver） |
| `public/` | 日志、存储、文档、测试资源 |
| `web/` | 前端工程（vue-pure-admin） |