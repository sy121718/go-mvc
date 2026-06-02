# CLAUDE.md

本文件描述当前仓库的实际开发约定。

## 项目概览

- 项目类型：Go + Vue 3 全栈 Web 项目
- 后端框架：Gin
- 前端框架：vue-pure-admin (Vue 3 + TypeScript + Vite)
- 运行模式：显式组件生命周期
- 目录结构：`cmd` / `config` / `internal` / `pkg` / `web` / `public`

## 常用命令

```bash
# 后端
go run cmd/main.go
go build -o app cmd/main.go
go test ./...
go test ./internal/module/backend/...

# 前端
cd web && pnpm dev
cd web && pnpm build
```

## 架构约定

### 1. 启动与关闭

统一入口：

- `config.Init()`
- `config.InitComponents()`
- `config.CloseComponents()`

组件不自行决定进程退出，组件只返回 `error`。
组件自己的严格配置校验在各自 `pkg.Init()` 内部完成，`config` 不手工点名外调校验函数。

### 2. 目录职责

- `cmd/`：启动入口
- `config/`：配置读取、默认值、组件注册与关闭编排
- `internal/routers/`：主路由聚合
- `internal/middleware/`：默认中间件与系统中间件
- `internal/module/`：业务模块
- `internal/task/`：任务注册与调度
- `pkg/`：可复用基础组件 facade
- `public/`：日志、存储、测试资源等公共目录

### 3. pkg 组件方向

`pkg` 当前以 facade 形式对外暴露统一入口，具体实现放在 provider/driver 子目录。

已有代表：

- `pkg/auth` — JWT 认证 + Redis 用户会话（封禁、在线心跳、踢人）
- `pkg/casbin` — RBAC 权限引擎（策略存 DB，启动加载到内存）
- `pkg/cache` — Redis 缓存 facade
- `pkg/database` — 数据库 facade
- `pkg/queue` — 任务队列 facade
- `pkg/upload` — 文件上传 facade
- `pkg/lock` — 分布式锁 facade

## 认证与鉴权架构

### 三层分工

```
JWT    → 认证（验签 + 解析 user_id，写入 gin.Context）
Casbin → 鉴权（内存 enforcer，Enforce(user_id, path, method)）
Redis  → 运行时状态（封禁标记、在线心跳、用户信息缓存）
```

### 认证流程

```
登录 → 验密 → JWT(含 user_id) → 写 Redis session → 返回 token
请求 → JWTAuthMiddleware → ParseToken → c.Set("user_id") → c.Next()
     → RefreshOnline(写心跳) → X-New-Token(自动续期)
```

### 鉴权流程

```
CasbinMiddleware:
  c.Get("user_id") → casbin.GetEnforcer().Enforce(user_id, path, method)
  → Casbin 内存 enforcer（启动时从 sys_casbin_rule 加载）
  → g 映射自动合并用户直接权限 + 角色继承权限
```

### 权限设计（菜单 = 权限的可视化）

- `sys_menus`: type 1=目录 2=菜单 3=按钮 4=外链。type=2,3 必含 permission_code
- 给角色/用户分配菜单 → 收集 type=2,3 的 permission_code → 去重 → 写入 sys_casbin_rule
- Casbin 策略: `p, target_id, path, method, code`
- 不需要 sys_menu_permission、sys_user_permission_exception

### Redis 用途

| Key | 用途 |
|-----|------|
| `user:blocked:{id}` | 封禁标记（中间件每次请求检查） |
| `online:{id}` | 在线心跳（5 分钟 TTL，管理后台展示） |
| `user:session:{id}` | 用户信息（Profile 接口优先查 Redis） |

---

## 配置约定

- 默认值统一放在 `config/config.go`
- `pkg` 自己定义并解析自己的配置结构
- `pkg` 不导入 `config` 包，避免循环依赖
- `config` 只负责读取配置和调度组件初始化，不承载具体 `pkg` 配置校验逻辑

## 路由与中间件

### 路由

- 只使用 `GET` 和 `POST`
- 主路由聚合在 `internal/routers/routes.go`
- 模块路由各自维护在模块自己的 `router/` 目录
- 当前默认健康检查路由：
  - `GET /livez`
  - `GET /readyz`

### 默认中间件

当前已实现的中间件能力：

- JWTAuthMiddleware — JWT 验签、封禁检查、自动续期、在线心跳
- CasbinMiddleware — RBAC 鉴权（需 JWT 前置）
- 安全响应头、请求体大小限制、固定窗口限流、Recovery

## 响应约定

当前实现已经调整为：

- `pkg` 和系统相关的包（config、cmd、middleware）不使用国际化翻译
- `pkg` 默认错误提示直接写最终中文，或直接返回原始系统错误
- 业务模块（internal/module）可以直接使用 `pkg/i18n` 读取多语文案
- `pkg/response` 使用数字状态码
- `pkg/response` 不再维护字符串错误码
- `pkg/response` 不再做国际化翻译
- 调用点直接写最终中文提示

标准结构：

```go
type Response struct {
Code    int         `json:"code"`
 Message string      `json:"message"`
Data    interface{} `json:"data,omitempty"`
}
```

推荐用法：

```go
response.Success(c, data)
response.SuccessWithMessage(c, "保存成功", data)
response.ErrorWithMessage(c, 400, "请求参数错误")
response.ErrorWithMessage(c, 401, "未登录或登录已过期")
response.ErrorWithMessage(c, 403, "无权限访问")
response.ErrorWithMessage(c, 404, "请求的资源不存在")
```

不要再写：

```go
response.Error(c, enums.ErrSystemError)
response.ErrorWithMessage(c, enums.ErrInvalidParams, "请求已过期")
```

## i18n 约定

`pkg/i18n` 仍然保留，业务模块（internal/module）可以直接使用：

- 读取数据库中的多语言文本（UI 文案、字典、业务文案）
- 查询指定语言下的文本内容

以下场景**仍然不使用** i18n（直接写中文提示或原始错误）：

- `pkg` 默认错误返回
- 系统中间件默认错误返回
- `pkg/response` 文案翻译
- config / cmd 等系统层初始化错误返回

## 错误处理约定

### 保留直接中文提示的场景

- 参数缺失
- 配置缺失
- 配置无效
- 状态不满足
- provider / driver 不支持
- 文件为空、扩展名不允许等规则校验

### 不要过度包装的场景

底层系统错误尽量直接返回原始 `err`，不要在 `pkg` 里重复翻译。

不建议：

```go
return fmt.Errorf("创建日志目录失败: %w", err)
```

更倾向：

```go
return err
```

另外：

- `pkg` 初始化失败，直接返回 `error`
- 启动链路收到错误后直接停止进程
- 不在外部重复做一层 `pkg` 配置校验中转

## 测试约定

- 默认跑现有测试，不新增额外测试框架
- 接口完成后优先维护 feature/usecase 测试
- 只有定位复杂函数问题时，再补定向单测

## git 提交约定

每次 git commit 必须在 message 中用中文注明：

- 改了哪个文件（路径）
- 改了什么（一句话简述）

格式示例：

```text
feat: 添加 Casbin 策略的 code 字段

pkg/casbin/casbin.go | p = sub, obj, act → p = sub, obj, act, code
```

## 模块开发

已实现的业务模块：

| 模块 | 路由前缀 | 说明 |
|------|---------|------|
| `internal/module/backend/admin` | `/api/admin` | 管理员 CRUD、登录、Profile |
| `internal/module/common/captcha` | `/api/captcha` | 图形验证码 |

涉及业务模块开发时，先看：

- [internal/module/CLAUDE.md](./internal/module/CLAUDE.md)

涉及 `pkg` 组件开发时，先看：

- [pkg/CLAUDE.md](./pkg/CLAUDE.md)
