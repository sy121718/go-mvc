# pkg 组件包说明

本目录包含项目的公共组件，提供可复用的基础能力模块，统一由 `config/config.go` 进行生命周期编排。

## 目录结构

```text
pkg/
├── auth/         # JWT 认证组件
├── cache/        # 缓存组件（根入口 + provider 实现）
├── casbin/       # 权限组件
├── crypto/       # 加密签名组件
├── database/     # 数据库组件（根入口 + driver 实现）
├── enums/        # 全局常量（含错误码）
├── errors/       # 通用错误定义
├── i18n/         # 多语言配置中心（数据库驱动）
├── logger/       # 日志组件
├── queue/        # 异步任务队列（根入口 + provider 实现）
├── response/     # 统一响应格式
├── upload/       # 上传组件（多 provider）
├── utils/        # 通用工具函数
└── validate/     # 通用校验能力
```

## 组件列表

### 1. auth - JWT 认证

**主要功能：**

- Token 生成、解析、刷新
- 用户身份验证
- 启动时完成配置装配

**使用示例：**

```go
// 生成 Token
token, err := auth.GenerateToken(userID, username)

// 解析 Token
claims, err := auth.ParseToken(token)
```

---

### 2. cache - 缓存组件

**主要功能：**

- 根入口位于 `pkg/cache/cache.go`
- 具体实现位于 `pkg/cache/provider/`
- 当前默认实现为 Redis provider
- 由 `redis.enabled` 控制是否启用
- 通过 `redis.provider` 选择具体 provider

**使用示例：**

```go
rdb := cache.GetRedis()
rdb.Set(ctx, "key", "value", time.Hour)
rdb.Get(ctx, "key")
```

---

### 3. casbin - 权限管理

**主要功能：**

- 基于数据库策略表加载权限规则
- 提供运行期权限校验入口
- 启用时在服务启动阶段完成初始化

**使用示例：**

```go
enforcer := casbin.GetEnforcer()
ok, err := enforcer.Enforce(sub, obj, act)
```

---

### 4. crypto - 加密签名

**主要功能：**

- API 签名验证
- HMAC-SHA256 加密

**使用示例：**

```go
sign := crypto.GenerateSign(params, timestamp)
crypto.VerifySign(params, sign, timestamp)
```

---

### 5. database - 数据库组件

**主要功能：**

- 根入口位于 `pkg/database/database.go`
- 具体驱动实现位于 `pkg/database/driver/`
- 连接管理
- 连接池配置
- GORM 集成
- 启动时完成连接检查
- 通过 `database.driver` 选择具体数据库驱动

**使用示例：**

```go
db := database.GetDB()
db.Find(&users)
```

---

### 6. enums / errors - 错误契约与通用错误

**职责范围：**

- `pkg/enums/errors.go` 作为全局错误码唯一来源。
- `pkg/errors` 保留 pkg 层通用错误能力，不作为业务错误码常量定义位置。
- 业务提示文案与 HTTP 状态码统一走 `pkg/i18n` 对应的 `sys_i18n` 字典。

---

### 7. i18n - 多语言配置中心

**主要功能：**

- 从数据库 `sys_i18n` 预热多语言缓存
- 统一查询消息文本与 HTTP 状态码
- 支持自动刷新与显式停止
- 默认语言优先读取 `config.yaml` 的 `i18n.default_lang`
- pkg 内部保留默认语言兜底

**使用示例：**

```go
import (
    "go-mvc/pkg/enums"
    "go-mvc/pkg/i18n"
)

code := enums.ErrUploadConfigMissing
result := i18n.Get(code, "zh-CN")   // 完整结构
text := i18n.GetText(code, "zh-CN") // 仅文案
httpCode := i18n.GetHttpCode(code)  // 仅状态码
```

---

### 8. queue - 异步队列

**主要功能：**

- 根入口位于 `pkg/queue/queue.go`
- 具体实现位于 `pkg/queue/provider/`
- 异步任务投递
- 延迟任务
- 可选 worker 启动
- 通过 `queue.provider` 选择具体 provider

**使用示例：**

```go
queue.Register("email:send", HandleEmailSend)
err := queue.Enqueue("email:send", payload)
err = queue.EnqueueIn("email:send", time.Minute, payload)
```

---

### 9. response - 统一响应

**主要功能：**

- 标准化 API 响应格式
- 统一依赖 i18n 输出文案与状态码

**使用示例：**

```go
import "go-mvc/pkg/enums"

response.Success(c, user)
response.SuccessWithMessage(c, "msg_operation_success", user)
response.Error(c, enums.ErrSystemError)
```

---

### 10. validate - 通用校验

**主要功能：**

- 提供通用校验辅助能力
- 供业务层按需复用

---

## 生命周期规范

### 1. 统一启动入口

所有基础组件统一由 `config.InitComponents()` 编排初始化，而不是在 `main.go` 中逐个直接调用。

```go
if err := config.Init("config.yaml"); err != nil {
    return err
}

if err := config.InitComponents(); err != nil {
    return err
}
```

### 2. 统一关闭入口

所有需要释放的组件统一由 `config.CloseComponents()` 逆序关闭。

```go
if err := config.CloseComponents(); err != nil {
    log.Printf("组件关闭失败: %v", err)
}
```

### 3. 组件初始化规则

- **强制启动初始化**：`database`、`auth`、`i18n`
- **按配置启用，启用后启动初始化**：`casbin`、`cache`、`queue`
- 不再以 `lazy_init` 作为通用设计原则
- 组件内部只负责初始化自身，不决定退出进程

### 4. 配置读取规则

- 配置默认值统一放在 `config/config.go`
- 各 pkg 自己定义 `Config` 结构体并解析对应配置段
- pkg 不导入 `config` 包，避免循环依赖

---

## 使用流程

### 1. 在启动链中初始化

```go
if err := config.Init("config.yaml"); err != nil {
    return err
}

serverCfg, err := config.GetServer()
if err != nil {
    return err
}

gin.SetMode(serverCfg.Mode)

if err := config.InitComponents(); err != nil {
    return err
}
```

### 2. 在业务代码中使用

```go
db := database.GetDB()
rdb := cache.GetRedis()
token, err := auth.GenerateToken(userID, username)
response.Success(c, data)
```

---

## 扩展新组件

### 1. 创建组件目录

新增组件放在 `pkg/<name>/` 下，并在组件内部定义自己的配置结构与初始化函数。

### 2. 设计初始化接口

推荐形式：

```go
func Init(v *viper.Viper) error {
    var cfg Config
    if err := v.UnmarshalKey("newpkg", &cfg); err != nil {
        return err
    }
    return nil
}
```

### 3. 接入生命周期编排

- 在 `config/config.go` 中设置默认值
- 在 `config.InitComponents()` 中按依赖顺序接入初始化
- 如需释放资源，在 `config.CloseComponents()` 中接入关闭逻辑

---

## 注意事项

1. **配置结构体定义在各自 pkg 中**
2. **避免循环导入：pkg 不导入 config 包**
3. **pkg 只返回 error，不决定退出进程**
4. **生命周期由 config 统一编排**
5. **`i18n` 继续以数据库为唯一数据源**

---
