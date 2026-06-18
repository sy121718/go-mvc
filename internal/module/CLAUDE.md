# internal/module 开发规范

`internal/module/` 存放业务模块代码。

## 目录结构

```text
module_name/
├── contract/
│   ├── <module>_service.go
│   └── <dependency>_<role>.go
├── inbound/
│   ├── http/
│   │   ├── <module>_handle.go
│   │   ├── <module>_router.go
│   │   └── <module>_provider.go
│   ├── rpc/
│   ├── mq/
│   └── cron/
├── outbound/
│   ├── <dependency>/
│   │   └── <dependency>_client.go
│   ├── mq/
│   ├── sdk/
│   └── cache/
├── service/
│   ├── <module>_service.go
│   └── <module>_<action>.go
├── model/
│   └── <module>_model.go
├── dto/
│   ├── <module>_req.go
│   └── <module>_resp.go
└── enums/              # 必选目录，统一存放模块响应消息、业务错误消息、i18n key
    └── <module>_enums.go
```

## 核心关系

- `inbound` 承接外部调用
- `service` 实现本模块对外暴露契约，并调用本模块对外依赖契约
- `outbound` 实现本模块对外依赖契约
- `model` 只做本模块数据库访问
- `contract` 只放抽象
- `dto` 放请求 / 响应结构
- `enums` 必须存在，统一管理模块响应内容

## 命名规则

- 对外暴露契约：`<module>_service.go`
- 对外依赖契约：`<dependency>_<role>.go`
- `role` 统一用：`provider` / `reader` / `writer` / `publisher`
- `inbound/http/`：`<module>_handle.go`、`<module>_router.go`、`<module>_provider.go`
- `service/`：`<module>_service.go` + `<module>_<action>.go`
- `model/`：`<module>_model.go`
- `dto/`：`<module>_req.go`、`<module>_resp.go`
- `enums/`：`<module>_enums.go`

## contract

- 两类契约分文件写，不混在一个 `*_contract.go`
- `service` 只依赖本模块 `contract/`
- `contract/` 不用 `client` 后缀
- 历史 `*_contract.go` 可兼容，新代码按新命名执行

## inbound/http

- `handle`：绑定参数、基础校验、调用 `service`、输出响应
- `router`：只注册路由
- `provider`：做依赖装配
- 返回给前端的响应消息统一取 `enums`

## service

- `xxx_service.go` 只放 `Service` / `Deps` / `NewService()`
- 业务用例拆到 `xxx_<action>.go`
- 返回 `error`
- 业务错误消息统一取 `enums`
- 不直接承载 RPC / HTTP / MQ 实现细节

## model

- 只做数据库访问
- 不放业务规则

## dto

- `*_req.go` 给 `inbound` 绑定
- `*_resp.go` 给 `service` 返回
- 数据流：`inbound -> service -> inbound`

## enums

- `enums/` 是必须目录
- 所有响应内容都走模块 `enums`
- 包括：成功消息、参数错误消息、未授权消息、业务错误消息
- `handle` 和 `service` 不直接硬编码响应文案
- 未接好 `i18n` 时，`ErrXxx` / `MsgXxx` 直接等于中文常量
- 接好 `i18n` 后，再把枚举值切到 i18n key 或 i18n 取值
- 对外调用方只认模块 `enums`，不要绕过 `enums` 直接取文案

## 响应

- 统一走 `pkg/response`
- 优先使用：`response.Success`、`response.SuccessWithMessage`、`response.ErrorWithMessage`
- 传给 `response` 的消息统一来自模块 `enums`

## 路由

- 只用 `GET` / `POST`
- `GET` 查询
- `POST` 用于新增、修改、删除、状态变化
