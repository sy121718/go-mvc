# Go MVC

一个前后端分离的全栈 Web 项目。

- 后端：Go + Gin
- 前端：Vue 3 + TypeScript + Vite + vue-pure-admin
- 目标：把启动编排、基础组件、业务模块、测试资源和前端页面按职责拆开，尽量让运行链路清楚、模块边界清楚

## 当前状态

当前仓库已经落地的核心结构：

- 后端统一入口在 `cmd/main.go`
- 配置与组件生命周期统一由 `config/` 编排
- 全局中间件统一在 `internal/middleware/`
- 主路由聚合在 `internal/routers/routes.go`
- 业务模块集中在 `internal/module/`
- 通用基础组件集中在 `pkg/`
- 公共产物、迁移、测试资源集中在 `public/`
- 前端工程独立放在 `web/`

## 技术栈

### 后端

- Gin
- GORM
- Viper
- JWT
- Casbin
- Redis
- Zap / lumberjack
- Asynq（已接入组件，按配置启用）

### 前端

- Vue 3
- TypeScript
- Vite
- Element Plus
- Pinia
- Axios
- vue-i18n
- TailwindCSS + SCSS
- vue-pure-admin / pure-admin-thin

## 目录结构

```text
go-mvc/
├── cmd/                  # 启动入口
├── config/               # 配置读取、运行时组件注册、初始化与关闭编排
├── internal/
│   ├── middleware/       # 全局中间件装配 + 内建中间件实现
│   ├── module/           # 业务模块
│   ├── routers/          # 主路由聚合
│   └── task/             # 异步任务注册
├── pkg/                  # 通用基础组件（facade + provider/driver）
├── public/
│   ├── backup/           # 备份与初始化资源
│   ├── migrations/       # 数据迁移
│   └── test/             # 测试代码、测试支撑、fixtures
├── web/                  # 前端工程
├── config.yaml.example   # 配置样例
├── CLAUDE.md             # 项目级开发约定
└── README.md
```

## 后端启动链路

当前后端启动过程是显式编排的，不让组件自行决定进程退出。

运行顺序：

1. `config.Init("config.yaml")`
2. `config.GetServer()`
3. `config.InitComponents()`
4. `middleware.Setup(router)`
5. `routers.SetupRoutes(router, config.ValidateReady)`
6. `http.Server` 启动
7. 收到退出信号后，先关 HTTP，再执行 `config.CloseComponents()`

对应职责：

- `config.Init()`：读取配置文件，初始化全局配置对象
- `config.InitComponents()`：按顺序初始化组件
- `config.ValidateReady()`：给 `/readyz` 健康检查使用
- `config.CloseComponents()`：逆序关闭已初始化组件

## 运行时组件

当前组件注册集中在 `config/register.go`，主要包括：

- `logger`
- `validate`
- `database`
- `i18n`
- `casbin`（按配置启用）
- `cache`（按配置启用）
- `auth`
- `upload`（按配置启用）
- `queue`（按配置启用）
- `captcha`

组件约定：

- 配置校验在各自 `pkg.Init()` 内完成
- 组件只返回 `error`，不自己退出进程
- 关闭顺序与初始化顺序相反

## 默认中间件链

`middleware.Setup()` 当前会挂载这条全局链路：

1. `Recovery`
2. `CORS`
3. `SecurityHeaders`
4. `RequestBodyLimit`
5. `RequestRateLimit`（按配置启用）
6. `RequestLogCapture`（按配置启用）

说明：

- `JWTAuthMiddleware` 已实现，但不是全局中间件，而是按路由组挂载
- `CasbinMiddleware` 已实现，但当前默认 `admin` 路由还没有接入
- `SignatureMiddleware` 已实现，但当前没有默认挂到主路由链路上

所以仓库里同时存在：

- 已全局生效的基础安全与治理能力
- 已实现但尚未接入到默认业务路由的能力

## 路由现状

主路由聚合在 `internal/routers/routes.go`。

当前默认注册的入口有：

- `GET /livez`
- `GET /readyz`
- `GET /api/captcha`
- `POST /api/admin/login`
- `GET /api/admin/list`
- `GET /api/admin/detail`
- `POST /api/admin/create`
- `POST /api/admin/edit`
- `GET /api/admin/profile`

路由风格约定：

- 只使用 `GET` / `POST`
- 不使用 RESTful 路径参数
- 查询参数统一走 query 或 JSON body

## 认证、鉴权与会话现状

### 认证

当前登录主链路在 `internal/module/backend/admin/service/admin_login.go`。

登录过程：

1. 校验验证码
2. 按用户名/邮箱查管理员
3. 检查锁定状态、禁用状态
4. 校验密码
5. 登录失败累计计数，连续 5 次失败封禁 30 分钟
6. 登录成功后生成 JWT
7. 写入用户会话
8. 刷新在线状态

### Token 下发与续期

当前后端登录成功后会同时返回两份 token：

- 响应头：`X-New-Token`
- 响应体：`data.accessToken`

自动续期逻辑：

- `JWTAuthMiddleware` 在请求结束后检查 token 剩余有效期
- 当剩余时间不超过 10 分钟时，重新签发 token
- 新 token 继续通过 `X-New-Token` 响应头返回

### 前端 token 存储现状

当前前端实现并不是“只放内存”。

实际行为是：

- Axios 请求前从 `@/utils/auth.ts` 读取 token
- token 会缓存到内存变量
- 同时也会落到本地存储
- 刷新页面后可以从本地存储恢复
- 收到 401 时会执行登出并清理 token

所以这里的真实状态是：

- 内存缓存：有
- 本地存储恢复：也有

### 鉴权

项目的设计约定是：

- JWT：负责认证
- Casbin：负责鉴权
- Redis：负责在线状态、封禁标记、用户会话缓存

但按当前默认路由现状看：

- `admin` 受保护接口已经接了 `JWTAuthMiddleware`
- `CasbinMiddleware` 已实现，但默认 `admin` 路由还未接入

如果后续要把权限链路补完整，应该是：

- 先 `JWTAuthMiddleware`
- 再 `CasbinMiddleware`
- 最后进入具体业务 handler

## 响应约定

统一响应结构：

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

当前约定：

- `pkg/response` 只负责输出结构
- 使用数字状态码
- 不做字符串错误码中转
- 不在 `pkg/response` 层做 i18n
- 业务模块的响应文案优先通过模块 `enums` 提供

## 业务模块分层约定

目标分层在 `internal/module/CLAUDE.md` 中已经明确：

```text
module_name/
├── contract/
├── inbound/
│   └── http/
├── outbound/
├── service/
├── model/
├── dto/
└── enums/
```

职责关系：

- `inbound`：承接外部输入
- `service`：承接业务用例
- `outbound`：实现外部依赖访问
- `model`：只做数据库访问
- `contract`：放抽象契约
- `dto`：请求 / 响应结构
- `enums`：统一管理模块响应消息与业务错误消息

### 当前仓库现状补充

当前代码总体在向上述分层收敛，但仍有少量历史结构尚未完全统一：

- `backend/admin` 已基本接近目标分层，已完成简化：构造器直传参数（不用 Deps）、装配逻辑内联到 router
- `common/captcha` 仍是较老的 `handle/router` 结构
- `admin/contract/admin_contract.go` 已更新到新命名规范

README 以“现状 + 约定”同时说明，避免把“规划”误写成“已全部完成”。

## pkg 组件约定

`pkg/` 采用 facade + provider/driver 的组织方式。

当前目录主要包括：

- `auth`
- `cache`
- `casbin`
- `crypto`
- `database`
- `enums`
- `i18n`
- `lock`
- `logger`
- `queue`
- `response`
- `upload`
- `utils`
- `validate`
- `captcha`

核心原则：

- 根包对外暴露统一 API
- 组件自己解析自己的配置
- 组件自己在 `Init()` 做校验
- 底层系统错误优先返回原始 `err`
- 不在 `pkg` 内扩散业务语义

## public 目录现状

`public/` 当前主要承担三类职责：

- 运行产物与公共资源
- 数据迁移与初始化资源
- 测试代码与测试支撑

当前已存在的重点子目录：

- `public/backup/`
- `public/migrations/`
- `public/test/`

说明：

- 项目规则里约定了 `docs/`、`logs/`、`storage/` 的职责
- 当前仓库里已经明确存在的是 `backup`、`migrations`、`test`
- 文档描述时优先以实际目录为准，同时保留规则约定作为开发方向

## 测试现状

测试代码当前集中在 `public/test/`，已落地的内容包括：

- `feature/`：接口链路测试
- `pkg/`：若干基础组件功能测试
- `support/`：测试初始化与测试客户端支撑
- `fixtures/`：测试资源

常用命令：

```bash
go test ./...
go test ./public/test/...
go test ./public/test/feature -v
```

测试约定方向：

- 接口链路测试优先
- 复杂规则再补 unit
- 不启动生产进程，测试进程内完成装配

## 前端工程结构

前端工程位于 `web/`，基于 vue-pure-admin / pure-admin-thin。

当前重点目录：

```text
web/
├── src/api/            # 接口定义
├── src/components/     # ReXxx 公共组件
├── src/views/          # 页面模块
├── src/store/          # Pinia 状态管理
├── src/router/         # 路由
├── src/layout/         # 布局
├── src/utils/          # HTTP、鉴权、工具方法
├── src/plugins/        # 插件注册
├── src/style/          # 全局样式
├── locales/            # 国际化文案
└── types/              # 全局类型声明
```

前端开发约定：

- `.vue` 只负责组件拼装
- 页面业务逻辑尽量下沉到 `hook.tsx` / `.ts`
- 优先复用 `ReXxx` 组件
- 视图层直接复用 `src/api/` 导出的请求类型

### 前端 API 现状补充

规则文件里把 `src/api/user.ts` 作为目标结构之一，但当前仓库实际文件是：

- `src/api/admin.ts`
- `src/api/menu.ts`
- `src/api/mock.ts`
- `src/api/routes.ts`

其中登录、验证码、当前用户信息请求目前仍放在 `src/api/admin.ts`，还没有单独拆出 `user.ts`。

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/sy121718/go-mvc.git
cd go-mvc
```

### 2. 准备后端配置

```bash
cp config.yaml.example config.yaml
```

然后按你的环境修改至少这些配置：

- `database.*`
- `redis.*`
- `jwt.secret`
- `log.*`

### 3. 启动后端

```bash
go mod download
go run cmd/main.go
```

默认启动后可访问：

- `http://localhost:8080/livez`
- `http://localhost:8080/readyz`

### 4. 启动前端

```bash
cd web
pnpm install
pnpm dev
```

### 5. 常用构建与检查

后端：

```bash
go build -o app cmd/main.go
go test ./...
```

前端：

```bash
cd web
pnpm typecheck
pnpm build
```

## 示例配置重点

`config.yaml.example` 当前已经覆盖这些配置段：

- `server`
- `database`
- `redis`
- `jwt`
- `casbin`
- `i18n`
- `queue`
- `log`
- `upload`

其中几个关键点：

- `server.request_body_limit`：普通请求体限制
- `server.upload_body_limit`：上传请求体限制
- `server.rate_limit_enabled`：是否启用固定窗口限流
- `queue.enabled`：是否启用任务队列
- `redis.enabled`：是否启用缓存组件
- `casbin.enabled`：是否启用 Casbin 组件

## 已知现状差异

为了避免文档误导，这里把几个容易想当然的点单独写清楚：

- token 当前不是只存在内存里，也会写入本地存储
- `CasbinMiddleware` 已实现，但默认 `admin` 路由未接入
- `SignatureMiddleware` 已实现，但默认主路由未接入
- `public/test` 的理想分层规范已经写好，但当前实际目录还处于逐步收敛阶段
- `web/src/api/user.ts` 还未落地，登录相关接口目前在 `web/src/api/admin.ts`
- `internal/module` 的新分层规范已明确，但仓库中仍保留少量旧结构

## 更多规则

更细的开发约定见：

- 项目总规则：`CLAUDE.md`
- 业务模块规则：`internal/module/CLAUDE.md`
- `pkg` 组件规则：`pkg/CLAUDE.md`
- `public` 目录规则：`public/CLAUDE.md`
- 测试规则：`public/test/CLAUDE.md`
- 前端规则：`web/CLAUDE.md`