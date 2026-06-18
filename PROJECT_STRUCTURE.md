# Go-MVC 项目目录结构

## 后端（已实现）

```text
go-mvc/
├── CLAUDE.md               # 项目级开发约定
├── config.yaml.example
├── LICENSE
├── PROJECT_STRUCTURE.md
├── README.md
├── cmd/
│   └── main.go
├── config/
│   ├── config.go
│   ├── register.go
│   └── runtime.go
├── internal/
│   ├── middleware/
│   │   ├── middleware.go
│   │   └── builtin/
│   │       ├── auth.go
│   │       ├── body_limit.go
│   │       ├── casbin.go
│   │       ├── cors.go
│   │       ├── rate_limit.go
│   │       ├── recovery.go
│   │       ├── request_log.go
│   │       ├── security_headers.go
│   │       └── signature.go
│   ├── module/
│   │   ├── CLAUDE.md          # 模块开发规范
│   │   ├── backend/
│   │   │   └── admin/
│   │   │       ├── contract/
│   │   │       │   └── admin_contract.go
│   │   │       ├── dto/
│   │   │       │   ├── admin_req.go
│   │   │       │   └── admin_resp.go
│   │   │       ├── enums/
│   │   │       │   └── admin_enums.go
│   │   │       ├── inbound/
│   │   │       │   └── http/
│   │   │       │       ├── admin_handle.go
│   │   │       │       └── admin_router.go
│   │   │       ├── model/
│   │   │       │   └── admin_model.go
│   │   │       ├── outbound/
│   │   │       │   ├── cache/
│   │   │       │   │   └── cache_client.go
│   │   │       │   ├── log/
│   │   │       │   │   └── log_client.go
│   │   │       │   └── permission/
│   │   │       │       └── permission_client.go
│   │   │       └── service/
│   │   │           ├── admin_create.go
│   │   │           ├── admin_delete.go
│   │   │           ├── admin_detail.go
│   │   │           ├── admin_edit.go
│   │   │           ├── admin_list.go
│   │   │           ├── admin_login.go
│   │   │           ├── admin_menu.go
│   │   │           ├── admin_profile.go
│   │   │           ├── admin_role.go
│   │   │           └── admin_service.go
│   │   └── common/
│   │       └── captcha/
│   │           ├── handle/
│   │           │   └── captcha_handle.go
│   │           └── router/
│   │               └── captcha_router.go
│   ├── routers/
│   │   └── routes.go
│   └── task/
│       ├── email.go
│       ├── order.go
│       └── register.go
├── pkg/
│   ├── auth/
│   │   ├── jwt.go
│   │   └── session.go
│   ├── cache/
│   │   ├── cache.go
│   │   └── provider/
│   │       ├── provider.go
│   │       └── redis.go
│   ├── captcha/
│   │   └── captcha.go
│   ├── casbin/
│   │   └── casbin.go
│   ├── crypto/
│   │   ├── hash.go
│   │   └── sign.go
│   ├── database/
│   │   ├── database.go
│   │   ├── gorm_logger.go
│   │   ├── utils.go
│   │   └── driver/
│   │       ├── config.go
│   │       ├── mysql.go
│   │       ├── postgres.go
│   │       ├── sqlite.go
│   │       └── sqlserver.go
│   ├── enums/
│   │   ├── errors.go
│   │   └── messages.go
│   ├── i18n/
│   │   ├── cache.go
│   │   ├── i18n.go
│   │   └── loader.go
│   ├── lock/
│   │   ├── local.go
│   │   ├── lock.go
│   │   └── redis.go
│   ├── logger/
│   │   ├── api.go
│   │   ├── core.go
│   │   ├── entry.go
│   │   ├── manager.go
│   │   └── types.go
│   ├── queue/
│   │   ├── queue.go
│   │   └── provider/
│   │       ├── asynq.go
│   │       └── provider.go
│   ├── response/
│   │   └── response.go
│   ├── upload/
│   │   ├── upload.go
│   │   └── provider/
│   │       ├── local.go
│   │       ├── provider.go
│   │       └── qiniu.go
│   ├── utils/
│   │   └── port.go
│   └── validate/
│       ├── register.go
│       ├── unique.go
│       ├── validate.go
│       └── provider/
│           └── rule_email.go
├── public/
│   ├── backup/
│   │   ├── database/
│   │   ├── demo/
│   │   └── json/
│   ├── migrations/
│   │   ├── 001_email_template.go
│   │   ├── 002_email_send_record.go
│   │   ├── 003_email_send_recipient.go
│   │   ├── 004_ip_blacklist.go
│   │   ├── 005_notice.go
│   │   ├── 006_notice_read.go
│   │   ├── 007_notice_target.go
│   │   ├── 008_sys_admin.go
│   │   ├── 009_sys_admin_social.go
│   │   ├── 010_sys_file_category.go
│   │   ├── 011_sys_attachment.go
│   │   ├── 012_sys_casbin_rule.go
│   │   ├── 013_sys_config.go
│   │   ├── 014_sys_cron_job.go
│   │   ├── 015_sys_i18n.go
│   │   ├── 016_sys_logs.go
│   │   ├── 017_sys_menus.go
│   │   ├── 018_sys_rule.go
│   │   ├── 019_sys_rule_assignment.go
│   │   ├── 020_seed_admin.go
│   │   ├── 021_sys_role.go
│   │   ├── 022_sys_permission.go
│   │   └── migrator.go
│   └── test/
│       ├── feature/
│       │   ├── body_limit_test.go
│       │   └── health_test.go
│       ├── fixtures/
│       ├── pkg/
│       │   ├── auth/
│       │   ├── cache/
│       │   ├── casbin/
│       │   ├── database/
│       │   ├── i18n/
│       │   ├── logger/
│       │   ├── queue/
│       │   ├── upload/
│       │   └── validate/
│       └── support/
│           ├── test_bootstrap.go
│           └── test_client.go
└── web/
    ├── .husky/
    ├── build/
    ├── locales/
    ├── mock/
    ├── public/
    ├── src/
    │   ├── api/
    │   ├── assets/
    │   │   ├── iconfont/
    │   │   ├── login/
    │   │   ├── status/
    │   │   ├── svg/
    │   │   └── table-bar/
    │   ├── components/
    │   ├── config/
    │   ├── directives/
    │   ├── layout/
    │   │   ├── components/
    │   │   │   ├── lay-content/
    │   │   │   ├── lay-footer/
    │   │   │   ├── lay-navbar/
    │   │   │   ├── lay-notice/
    │   │   │   ├── lay-panel/
    │   │   │   ├── lay-search/
    │   │   │   ├── lay-setting/
    │   │   │   ├── lay-sidebar/
    │   │   │   └── lay-tag/
    │   │   └── hooks/
    │   ├── plugins/
    │   ├── router/
    │   │   └── modules/
    │   ├── store/
    │   │   └── modules/
    │   ├── style/
    │   ├── utils/
    │   │   └── http/
    │   └── views/
    └── types/
```

## 后端（规划中）

```text
internal/module/backend/
├── user/                  # 前端用户管理
│   ├── contract/
│   ├── dto/
│   ├── enums/
│   ├── inbound/http/
│   ├── model/
│   ├── outbound/
│   └── service/
├── menu/                  # 菜单管理
│   ├── contract/
│   ├── dto/
│   ├── enums/
│   ├── inbound/http/
│   ├── model/
│   └── service/
├── role/                  # 角色管理
│   ├── contract/
│   ├── dto/
│   ├── enums/
│   ├── inbound/http/
│   ├── model/
│   └── service/
├── notice/                # 系统通知
│   ├── contract/
│   ├── dto/
│   ├── enums/
│   ├── inbound/http/
│   ├── model/
│   └── service/
├── file/                  # 文件管理
│   ├── contract/
│   ├── dto/
│   ├── enums/
│   ├── inbound/http/
│   ├── model/
│   └── service/
├── log/                   # 操作日志
│   ├── contract/
│   ├── dto/
│   ├── enums/
│   ├── inbound/http/
│   ├── model/
│   └── service/
├── config/                # 系统配置管理
│   ├── contract/
│   ├── dto/
│   ├── enums/
│   ├── inbound/http/
│   ├── model/
│   └── service/
└── task/                  # 定时任务管理
    ├── contract/
    ├── dto/
    ├── enums/
    ├── inbound/http/
    ├── model/
    └── service/

pkg/
├── sms/                   # 短信发送
│   ├── sms.go
│   └── provider/
├── payment/               # 支付
│   ├── payment.go
│   └── provider/
├── websocket/             # WebSocket
│   └── websocket.go
└── search/                # 全文搜索
    ├── search.go
    └── provider/

public/test/backend/
├── admin/
│   ├── feature/
│   └── unit/
└── user/
    ├── feature/
    └── unit/
```

---

## 目录职责说明

### 后端顶层

| 目录 | 职责 |
|------|------|
| `CLAUDE.md` | 项目级开发约定与规范 |
| `PROJECT_STRUCTURE.md` | 项目目录结构总览（本文件） |
| `README.md` | 项目说明与技术栈 |
| `config.yaml.example` | 配置文件模板 |
| `cmd/` | 应用启动入口，仅含 `main.go` |
| `config/` | 配置读取、组件注册表、运行时生命周期编排 |
| `internal/` | 内部实现，不对外暴露 |
| `pkg/` | 可复用基础组件库（facade + provider/driver） |
| `public/` | 公共资源、数据迁移、测试代码与支撑 |

### internal/ 子目录

| 目录 | 职责 |
|------|------|
| `middleware/` | 全局中间件装配入口 + 内建中间件实现 |
| `module/CLAUDE.md` | 业务模块开发规范（分层、命名、约定） |
| `module/backend/` | 后端管理类业务模块（admin / 待建模块） |
| `module/common/` | 通用模块（验证码等） |
| `routers/` | 主路由聚合入口 |
| `task/` | 异步任务注册与处理器 |

### 业务模块分层（每个模块内部）

| 子目录 | 职责 |
|------|------|
| `contract/` | 接口契约：对外暴露契约 + 对外依赖契约 |
| `inbound/http/` | HTTP 入口：`handle` 参数绑定与响应 + `router` 装配与路由注册 |
| `service/` | 业务逻辑实现 |
| `model/` | 数据库访问（GORM） |
| `dto/` | 请求/响应结构体 |
| `enums/` | 模块级消息常量与错误常量 |
| `outbound/` | 对外部依赖的调用封装（缓存 / 日志 / 权限等） |

### pkg/ 组件

| 目录 | 职责 |
|------|------|
| `auth/` | JWT 签发、验签、续期；用户会话管理 |
| `cache/` | 缓存 facade：统一 API，Redis 驱动 |
| `captcha/` | 图形验证码生成与校验 |
| `casbin/` | 权限引擎：内存 enforcer，Enforce API |
| `crypto/` | 密码哈希（bcrypt）、请求签名 |
| `database/` | 数据库 facade：多驱动支持（MySQL/PG/SQLite/SQLServer）、GORM 日志、工具函数 |
| `enums/` | 全局常量仓库（历史兼容，新模块走各自 enums） |
| `i18n/` | 后端国际化：多语言文案加载、缓存、查询 |
| `lock/` | 锁 facade：本地互斥锁、Redis 分布式锁 |
| `logger/` | 日志管理器：zap + lumberjack 滚动、多 writer |
| `queue/` | 消息队列 facade：Asynq 驱动 |
| `response/` | 统一 HTTP 响应结构（Code + Message + Data） |
| `upload/` | 文件上传 facade：本地存储、七牛云驱动 |
| `utils/` | 通用工具函数 |
| `validate/` | 数据校验 facade：自定义规则注册、唯一性校验 |

### public/ 子目录

| 目录 | 职责 |
|------|------|
| `backup/` | 数据库备份脚本、演示代码、初始化数据 JSON |
| `migrations/` | 数据库表结构迁移（22 张表 + migrator） |
| `test/feature/` | 接口链路测试 |
| `test/pkg/` | 基础组件功能测试 |
| `test/support/` | 测试环境初始化与测试客户端封装 |
| `test/fixtures/` | 静态测试资源 |

### 规划中模块

| 目录 | 计划职责 |
|------|---------|
| `module/backend/user/` | 前端用户管理 |
| `module/backend/menu/` | 菜单管理（CRUD + 权限码绑定 + 树形展示） |
| `module/backend/role/` | 角色管理（CRUD + 权限分配 + 用户关联） |
| `module/backend/notice/` | 系统通知（站内信、公告、已读管理） |
| `module/backend/file/` | 文件管理（上传、分类、附件关联） |
| `module/backend/log/` | 操作日志查询 |
| `module/backend/config/` | 系统动态配置管理 |
| `module/backend/task/` | 定时任务管理界面 |
| `pkg/sms/` | 短信发送 facade + 多驱动 |
| `pkg/payment/` | 支付 facade + 微信/支付宝驱动 |
| `pkg/websocket/` | WebSocket 连接管理 |
| `pkg/search/` | 全文搜索 facade + ES 驱动 |