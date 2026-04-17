# pkg 包审核待办清单

更新时间：2026-04-17

## 目标
- 记录本轮 `pkg` 代码审核发现的问题
- 按优先级拆分后续修复任务
- 提供明确验收标准，避免漏改

## P0（优先修复）

- [x] `pkg/validate/unique.go`：移除动态字段字符串拼接带来的注入风险
  - 涉及位置：`IsUnique` / `IsExists` / `IsUniqueExclude` / `IsUniqueExcludeField`
  - 处理结果：字段名改为白名单映射（业务字段 -> 真实列名），并使用 `gorm/clause` 构造条件

- [x] `pkg/validate/unique.go`：补齐数据库错误处理
  - 当前问题：`Count` 的 `Error` 未处理，DB 异常会被误判为“唯一/不存在”
  - 处理结果：函数统一返回 `(bool, error)`，调用方可显式处理异常

- [x] `pkg/logger/manager.go`：去掉 `ensureInited` 中被吞掉的初始化错误
  - 当前问题：`_ = Init(nil)` 静默失败
  - 处理结果：`ensureInited` 改为返回 `error`，初始化失败可观测

- [x] `pkg/logger/core.go`：目录创建失败时不要静默 fallback
  - 当前问题：`MkdirAll` 失败后继续回退且错误被吞掉
  - 处理结果：目录创建失败直接返回可诊断错误

- [x] `pkg/auth/jwt.go`：修复全局状态并发安全问题
  - 当前问题：`jwtSecret/jwtConfig/inited` 为包级可变状态，读写未加锁
  - 处理结果：增加 `sync.RWMutex`，初始化与读取走并发保护

- [x] `pkg/crypto/sign.go`：修复签名密钥可变全局变量并发风险
  - 当前问题：`secretKey` 可运行时修改，读写无同步保护
  - 处理结果：增加读写锁保护，并统一经快照读取密钥

## P1（次优先修复）

- [x] `pkg/database/database.go`：评估 `GetDB()` 直接 `panic` 的风险
  - 处理结果：移除 `panic` 风格实现，`GetDB()` 直接改为返回 `(*gorm.DB, error)`（不保留兼容入口）

- [x] `pkg/cache/provider/redis.go`：评估 `Client()` 直接 `panic` 的风险
  - 处理结果：移除 `panic` 风格实现，`Client()` 与 `GetRedis()` 直接改为返回 `(client, error)`（不保留兼容入口）

- [x] `pkg/logger/manager.go`：`Sync` 相关错误统一处理策略
  - 当前问题：多处 `Sync` 错误被忽略
  - 处理结果：统一收敛为错误聚合返回，重复实例关闭失败会写系统日志

- [x] `pkg/auth/jwt.go`：显式校验签名算法白名单
  - 处理结果：`ParseWithClaims` 增加 method 校验，仅允许 `HS256`

- [x] `pkg/response/response.go`：规范 `Accept-Language` 解析
  - 当前问题：直接使用整个 header 字符串，未处理多语言与 q 值
  - 处理结果：新增解析逻辑，按 `q` 权重选首选语言并标准化格式

- [x] `pkg/database/database.go`、`pkg/crypto/sign.go`：修复包注释格式
  - 当前问题：`// Package xxx /*` 非标准 GoDoc 形式
  - 处理结果：统一为标准包注释

## P2（优化项）

- [x] 为 `pkg/logger` 增加并发场景与异常路径测试
  - 处理结果：新增 `public/test/pkg/logger/logger_functional_test.go`
- [x] 为 `pkg/validate` 增加错误路径测试（DB 异常、非法字段）
  - 处理结果：新增 `public/test/pkg/validate/validate_functional_test.go`
- [x] 增加 `staticcheck` 与 `go test -race` 到常规检查流程
  - 处理结果：新增 `.github/workflows/quality-check.yml`
- [x] 逐步减少 `panic` 风格 API（改为显式返回 `error` 并迁移调用）
  - 处理结果：`config.InitComponents`、`pkg/i18n.LoadCache` 改为使用 `database.GetDB()`

## 验收标准

- [x] `go test ./...` 通过
- [x] `go vet ./...` 通过
- [x] `staticcheck ./...` 通过
- [x] 不再存在关键路径静默吞错
- [x] `validate` 不再接受未校验字段名直接拼接 SQL
- [ ] `go test -race ./...` 通过（当前本机为 `windows/386` 且缺少 `gcc`，已在 CI 工作流中启用）
