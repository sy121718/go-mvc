# 框架治理待办清单

更新时间：2026-04-21

## 目标

- 记录当前 Go MVC 基础框架的治理项
- 按优先级拆分运行时安全、结构化、性能与调用体验优化
- 作为后续框架迭代的统一 TODO，不混入具体业务需求

## 框架定位

- 这套框架是给自己长期使用的 Go MVC 多应用骨架
- 设计思路借鉴 ThinkPHP 多应用模式，但不复制 ThinkPHP 的魔法机制
- 保留 ThinkPHP 在目录组织、公共能力复用、开发顺手感上的优点
- 去掉隐式 hook、自动发现、静默注册、吞错等不透明机制
- 框架代码优先追求显式、可控、可排错，而不是少写几行代码

## 设计原则

- 借鉴 ThinkPHP 的组织效率，不接受 ThinkPHP 的魔法机制
- 所有初始化显式执行、所有注册显式调用、所有依赖显式传递、所有错误显式处理
- 禁止关键路径使用 `_` 跳过错误，禁止静默吞错
- 这是基础框架治理，不考虑旧接口兼容；发现结构不合理时，直接按目标结构收敛
- 锁分层处理，避免默认上 Redis 锁
- 先补运行时安全底线，再做性能治理
- 优先统一框架入口、生命周期、注册机制，再补业务层便利能力
- 默认行为显式化，避免静默 fallback 和隐式初始化
- 注册只能有一个入口，其他文件只负责辅助、启动、调用、编排、关闭
- 只有“系统启动要感知的包”才允许进入注册中心
- 纯工具包、纯业务包、DTO、model、helper 不进入系统注册中心
- 禁止通过 `init()` 承担关键运行时注册职责
- `config` 层按“读取、注册、调度”三块拆分：读取配置、注册清单、运行时调度
- 启动型协议统一必须落在 `pkg` 层，不能下沉到 `internal` 私有业务层
- `internal/task` 只负责模块间任务调度、自动任务、业务任务，不承担框架组件启动职责
- 已落地改动也要回头收敛，优先按最优结构改，不为小改保留过渡态代码

## P0（优先处理）

### 注册中心与启动骨架

- [x] 将 `config` 目录按“读取、注册、调度”职责拆解
  - 目标：形成 `config/config.go` 负责读取、`config/register.go` 负责注册、`config/runtime.go` 负责调度的清晰边界

- [x] 新增唯一注册入口 `config/register.go`
  - 目标：所有运行时组件与业务模块只能在一个地方注册，后续新增启动型包只改一个文件

- [ ] 明确注册准入原则：只允许“启动型包”进入注册中心
  - 目标：区分启动型、调用型、业务型三类包，避免普通工具包和业务包污染注册中心

- [x] 将组件注册与模块注册统一收口到 `config/register.go`
  - 目标：禁止在 `main`、`pkg`、`internal/module` 等其他位置分散维护运行时注册

- [ ] 禁止 `pkg` / `module` 通过 `init()` 承担关键运行时注册
  - 目标：消除隐式魔法，所有关键注册路径必须可见、可追踪

- [ ] 统一启动型 `pkg` 组件接入协议
  - 目标：所有启动型包统一暴露 `Init(cfg)` / `Close()`，初始化细节由包内部处理，注册中心只负责登记与顺序

- [ ] 明确启动型 `pkg` 与 `internal` 私有业务层的边界
  - 目标：组件启动协议、生命周期、配置解析留在 `pkg`；`internal/task` 仅保留任务调度、自动任务与业务任务

- [x] 清理 `config/register.go` 中的组件初始化胶水逻辑
  - 目标：注册中心尽量只保留 `Name`、`Enabled`、`Init`、`Close` 的直接绑定，不承担组件内部初始化细节

- [ ] 回看已落地改动并做二次收敛
  - 目标：把之前为了快速落地保留的包装层、过渡态结构、重复胶水代码继续收简，保持 Go 风格的直接、显式、少代码
  - 当前明确回看范围：
    1. `config/register.go`
       - 继续压缩组件注册清单，去掉无意义包装函数，尽量只保留 `Enabled`、`Init`、`Close` 的直接绑定
    2. `config/runtime.go`
       - 评估 `runtimeComponent`、`runtimeModule`、`initializedRegistry` 这一层抽象是否过重，能否进一步减少类型和状态管理代码量
       - 将 `/livez`、`/readyz`、`NoRoute`、`/api` 分组等 HTTP 路由装配细节迁出 `config`，回收到专门的路由层
    3. `pkg/queue/queue.go`
       - 回看 `registrations`、`registrationMu`、provider 重建后重放注册这套逻辑，判断能否在保证显式注册的前提下继续收简
    4. `internal/task/register.go`
       - 回看当前“显式任务注册”方案是否仍有绕路，确认是否还能再减少一层转发或集中度更高的注册写法
    5. `pkg/i18n/i18n.go`
       - 回看初始化和加锁流程，减少为了快速落地留下的双阶段加锁、配置分支和状态切换代码

- [x] 将 HTTP 路由装配细节迁出 `config`
  - 目标：把 `/livez`、`/readyz`、`NoRoute`、`/api` 分组等路由装配细节迁回专门的路由层，`config` 只保留模块注册清单和就绪判断

- [x] `pkg/queue`：收敛为统一组件协议样板
  - 目标：由 `pkg/queue` 自行处理初始化、是否启动 worker、关闭逻辑；禁止把队列组件启动职责下沉到 `internal/task`

- [x] `pkg/i18n`：收敛为统一组件协议样板
  - 目标：默认语言设置、自动刷新启动/停止收回 `pkg/i18n` 内部处理，减少注册中心胶水代码

- [x] `pkg/i18n`：收简重初始化与自动刷新状态切换逻辑
  - 目标：去掉双阶段加锁和旧刷新器残留问题，让重初始化按最新配置直接生效

- [x] 评估依赖型组件的统一协议接入方式
  - 目标：明确 `casbin` 这类依赖 DB 的组件如何在统一协议下处理前置依赖

- [x] 为启动型组件统一协议补功能测试
  - 目标：验证组件在统一协议下可初始化、可关闭、可按配置触发附加启动行为

- [x] `config/config.go`：将组件生命周期编排升级为统一组件接口与注册表
  - 目标：避免 `InitComponents` / `CloseComponents` 持续硬编码膨胀

- [x] 收缩 `config/config.go` 职责，仅保留配置读取、默认值、配置访问接口
  - 目标：将组件初始化、关闭、路由装配等运行时逻辑迁出 `config.go`

- [x] 新增 `config/runtime.go`，统一承接初始化、关闭、路由装配、健康检查、运行时状态管理
  - 目标：让 `register.go` 只管注册清单，`runtime.go` 只管执行调度

- [x] `internal/routers/routes.go`：建立模块路由统一注册机制
  - 目标：减少主路由中心文件手工 import 和手工调用

- [ ] 为业务模块增加统一入口文件，如 `module.go` / `bootstrap.go`
  - 目标：显式暴露模块路由注册、依赖初始化、ready 检查等能力

- [ ] 将启动流程显式拆分为 `LoadConfig -> ValidateConfig -> InitCritical -> InitOptional -> BuildRouter -> StartServer`
  - 目标：减少启动逻辑耦合，便于后续拆分关键与扩展启动阶段

- [ ] `cmd/main.go`：由 `gin.Default()` 改为 `gin.New()` + 显式挂载中间件
  - 目标：统一请求日志、恢复处理、安全中间件的注册顺序，避免默认中间件与自定义链路重复或失控

- [x] `pkg/queue`、`internal/task`：重构队列与任务调度边界
  - 目标：`pkg/queue` 负责队列组件协议与生命周期，`internal/task` 负责任务注册与任务业务，避免组件协议依赖私有调度层

### 运行时安全与服务保护

- [x] `cmd/main.go`：为 `http.Server` 增加 `ReadHeaderTimeout`、`ReadTimeout`、`WriteTimeout`、`IdleTimeout`
  - 目标：补齐最基础的连接保护与慢请求防护

- [x] `internal/middleware/auth.go`：错误响应后统一执行 `c.Abort()`
  - 目标：避免后续 handler 继续执行造成鉴权绕过

- [x] `internal/middleware/casbin.go`：错误响应后统一执行 `c.Abort()`
  - 目标：避免权限校验失败时请求链继续向下执行

- [x] `config/config.go`、`pkg/auth/jwt.go`：增加 release 模式 fail-fast 校验
  - 目标：禁止默认 JWT secret、默认数据库名、空关键配置在生产模式启动

- [x] `internal/routers/routes.go`：将 `/health` 拆分为 `/livez` 与 `/readyz`
  - 目标：区分存活检查与依赖就绪检查，支撑部署与故障切换

- [ ] 新增默认安全头中间件
  - 目标：统一输出 `X-Content-Type-Options`、`X-Frame-Options`、基础 `Content-Security-Policy`

- [ ] `internal/middleware/auth.go`：收敛 JWT 对外错误文案
  - 目标：统一返回“未认证”或“认证失效”，减少探测面

### 请求与输入防护

- [ ] 新增统一请求体大小限制中间件
  - 目标：区分普通 API 和上传接口的默认大小限制

- [ ] `pkg/upload`：增加最大大小限制、MIME 白名单、扩展名白名单的默认上传防护策略
  - 目标：把上传安全校验前置到框架层，而不是散落在业务层

- [ ] `internal/middleware/signature.go`：增加 `nonce` 去重机制
  - 目标：配合时间戳校验，防止 5 分钟窗口内请求重放

- [ ] 新增基础限流能力
  - 目标：先支持按 IP、按路由的默认限流，后续再扩展按用户限流

## P1（次优先处理）

### 结构化与组件边界

- [ ] `internal/middleware`：增加默认中间件聚合注册入口
  - 目标：安全头、限流、鉴权、签名、防重放等默认链路统一装配

- [ ] `pkg/response`、`pkg/enums`、`pkg/i18n`：统一消息码使用方式
  - 目标：禁止直接散写 `"msg_operation_success"` 这类字符串，统一走常量

- [ ] 将组件初始化拆为关键启动和扩展启动
  - 目标：数据库、JWT、日志作为关键依赖；i18n 自动刷新、queue worker、Casbin 预热作为扩展依赖

- [ ] 为各组件增加 `Ready()` 或等价就绪检查接口
  - 目标：避免将“调用过 Init”误当作“真实可用”

- [ ] `config/config.go`：增加测试态重置能力或独立 runtime 隔离
  - 目标：支持测试中切换配置与重复初始化，不依赖全局单例残留状态

- [ ] 统一启动型 `pkg` 的接口、错误语义和生命周期约定
  - 目标：让 `cache`、`queue`、`upload`、`i18n`、后续 `lock` 的使用体验保持一致

- [ ] 逐步从全局单例向 App Runtime 容器过渡
  - 目标：降低 `config`、`database`、`cache`、`auth`、`upload` 等包级状态对测试和扩展的限制

- [ ] `internal/module/backend/admin/model/admin_model.go`：移除 hook 中直接断言 context 值的 panic 风险
  - 目标：将 `tx.Statement.Context.Value(...).(bool)` 改为安全取值与默认拒绝策略，避免运行时崩溃

### 锁与并发治理

- [ ] 形成三层锁策略规范
  - 目标：进程内锁使用 `sync.Mutex` / `sync.RWMutex`；数据一致性优先唯一索引、事务、乐观锁；跨实例互斥再使用分布式锁

- [ ] 新建 `pkg/lock`
  - 目标：抽象统一锁接口，先支持 `local` 和 `redis` 两种 provider

- [ ] 明确 `pkg/lock` 适用边界
  - 目标：仅用于跨实例抢占、单实例任务执行、幂等提交、重复消费防护

- [ ] 热点并发控制优先引入 `singleflight`
  - 目标：用于缓存重建、热点查询合并，避免默认走 Redis 锁

### 调用方便与框架门面

- [ ] `pkg/cache`：补高层 helper，如 `GetOrLoad`、`Remember`、`SetJSON`
  - 目标：减少业务层直接操作底层 Redis client

- [ ] `pkg/queue`：补更高层任务门面
  - 目标：减少业务层手工拼 task type、payload、opts

- [ ] `pkg/upload`：补统一上传校验门面
  - 目标：把 provider 调度和安全校验分层，业务直接调用更稳定

- [ ] `pkg/logger`：取消隐式自动初始化
  - 目标：让初始化顺序问题在启动阶段显式暴露，不在运行期偷偷 fallback

- [ ] 统一组件错误分类
  - 目标：区分配置错误、初始化失败、未就绪、运行时错误，减少各包返回风格不一致

## P2（优化项）

### 性能治理

- [ ] `pkg/database/database.go`：增加可配置性能开关
  - 目标：支持 `PrepareStmt`、`SkipDefaultTransaction`、慢 SQL 阈值等能力

- [ ] 为数据库读写分离预留配置结构
  - 目标：先留好扩展位，避免未来改配置结构时破坏兼容

- [ ] `pkg/cache`：补缓存穿透、击穿、雪崩治理能力
  - 目标：包括空值缓存、TTL 抖动、`singleflight` 合并请求

- [ ] `pkg/logger`：增加 HTTP 和 SQL 日志采样开关与场景级别控制
  - 目标：降低高并发下全量日志对吞吐的影响

- [ ] `pkg/queue/provider/asynq.go`：将队列权重改为配置项
  - 目标：避免 `critical/default/low` 权重写死

### 可维护性与一致性

- [ ] 收敛默认值定义位置
  - 目标：减少 `config/config.go`、`pkg/database/driver/config.go`、`pkg/auth/jwt.go`、`pkg/cache/provider/redis.go`、`pkg/queue/provider/asynq.go` 之间的默认值漂移

- [ ] `pkg/i18n/cache.go`：避免将内部缓存 map 直接暴露给调用方
  - 目标：防止外部误改 `AllLangs` 污染内存缓存

- [ ] `pkg/utils/port.go`：将端口占用处理策略配置化
  - 目标：自动杀进程行为只允许在开发模式启用，并支持更清晰的策略切换

- [ ] 为框架常用入口补统一 facade 文档式注释与示例
  - 目标：提高团队成员对 `cache`、`queue`、`upload`、`lock` 等基础能力的调用一致性

## 执行顺序建议

- [ ] 第 1 批先做 P0 的注册中心与启动骨架
- [ ] 第 2 批做 P0 的运行时安全与请求防护
- [ ] 第 3 批做 P1 的结构化边界、锁策略与调用门面
- [ ] 第 4 批做 P2 的性能治理与一致性收敛
