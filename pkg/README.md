# pkg 组件包说明

本目录包含项目的公共组件，提供可复用的功能模块。

## 目录结构

```
pkg/
├── auth/         # JWT 认证组件
├── cache/        # Redis 缓存组件
├── crypto/       # 加密签名组件
├── database/     # MySQL 数据库组件
├── errors/       # 错误码定义
├── queue/        # 异步任务队列
└── response/     # 统一响应格式
```

## 组件列表

### 1. auth - JWT 认证

**主要功能：**
- Token 生成、解析、刷新
- 用户身份验证

**使用示例：**
```go
// 生成 Token
token, _ := auth.GenerateToken(userID, username)

// 解析 Token
claims, _ := auth.ParseToken(token)
```

---

### 2. cache - Redis 缓存

**主要功能：**
- Redis 连接管理
- 缓存操作

**使用示例：**
```go
rdb := cache.GetRedis(viper)
rdb.Set(ctx, "key", "value", time.Hour)
rdb.Get(ctx, "key")
```

---

### 3. crypto - 加密签名

**主要功能：**
- API 签名验证
- HMAC-SHA256 加密

**使用示例：**
```go
// 生成签名
sign := crypto.GenerateSign(params, timestamp)

// 验证签名
crypto.VerifySign(params, sign, timestamp)
```

---

### 4. database - MySQL 数据库

**主要功能：**
- 数据库连接管理
- 连接池配置
- GORM 集成

**使用示例：**
```go
db := database.GetDB(viper)
db.Find(&users)
```

---

### 5. errors - 错误码

**主要功能：**
- 标准化错误码
- 统一错误消息

**使用示例：**
```go
// 创建错误
err := errors.NewError(errors.UserNotFound)

// 自定义消息
err := errors.NewErrorWithMessage(errors.UserNotFound, "用户不存在")
```

---

### 6. queue - 异步队列

**主要功能：**
- 异步任务调度
- 定时任务
- 任务重试

**使用示例：**
```go
// 注册任务
queue.RegisterTask("email:send", HandleEmailTask)

// 投递任务
task := asynq.NewTask("email:send", payload)
queue.Enqueue(task)
```

---

### 7. response - 统一响应

**主要功能：**
- 标准化 API 响应格式
- 简化 Controller 代码

**使用示例：**
```go
// 成功响应
response.Success(c, user)

// 失败响应
response.Fail(c, "用户名已存在")

// 错误响应
response.Error(c, "服务器错误")
```

---

## 设计原则

### 1. 懒加载（Lazy Init）

所有组件支持懒加载，通过配置控制：

```yaml
database:
  lazy_init: false  # false=启动时连接，true=用时连接
```

### 2. 单例模式

使用 `sync.Once` 确保只初始化一次：

```go
var once sync.Once

once.Do(func() {
    // 只执行一次
})
```

### 3. 依赖注入

通过传递 `viper.Viper` 实例，pkg 自己解析配置：

```go
// main.go
v := config.GetViper()
database.InitDB(v)

// pkg/database/mysql.go
func InitDB(v *viper.Viper) error {
    var cfg Config
    v.UnmarshalKey("database", &cfg)
}
```

---

## 使用流程

### 1. 在 main.go 中初始化

```go
// 加载配置
config.Init("config.yaml")
v := config.GetViper()

// 初始化组件
database.InitDB(v)
cache.InitRedis(v)
auth.InitJWT(v)
```

### 2. 在业务代码中使用

```go
// 获取数据库连接
db := database.GetDB(viper)

// 获取 Redis 连接
rdb := cache.GetRedis(viper)

// 生成 Token
token, _ := auth.GenerateToken(userID, username)

// 返回响应
response.Success(c, user)
```

---

## 扩展新组件

### 1. 创建目录

```bash
mkdir pkg/newpkg
touch pkg/newpkg/newpkg.go
```

### 2. 定义配置结构体

```go
// pkg/newpkg/newpkg.go
package newpkg

type Config struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}
```

### 3. 实现初始化函数

```go
func InitNewPkg(v *viper.Viper) error {
    var cfg Config
    v.UnmarshalKey("newpkg", &cfg)
    // 初始化逻辑
    return nil
}
```

### 4. 在 config.yaml 添加配置

```yaml
newpkg:
  host: 127.0.0.1
  port: 8080
```

### 5. 在 main.go 调用

```go
newpkg.InitNewPkg(v)
```

---

## 注意事项

1. **配置结构体定义在各自的 pkg 中**
2. **避免循环导入**：pkg 不导入 config 包
3. **使用 viper 参数传递**：main.go 传递给 pkg
4. **并发安全**：使用 `sync.Once` 确保线程安全

---

## PHP 开发者对比

| 功能 | PHP (Laravel) | Go (本项目) |
|------|--------------|------------|
| 数据库 | `DB::table()` | `database.GetDB(v)` |
| 缓存 | `Cache::get()` | `cache.GetRedis(v).Get()` |
| 认证 | `JWTAuth::fromUser()` | `auth.GenerateToken()` |
| 响应 | `response()->json()` | `response.Success(c, data)` |
| 错误码 | `define('ERROR_CODE', 1000)` | `const ErrorCode = 1000` |
| 队列 | `dispatch(new Job())` | `queue.Enqueue(task)` |