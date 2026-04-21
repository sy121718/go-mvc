# Go MVC

一个基于 Gin 的 Go MVC 项目，当前采用显式组件生命周期、模块化目录结构和 `pkg` facade 入口。

## 当前重点

- 启动流程已经拆成显式步骤：`config.Init()`、`config.InitComponents()`、路由注册、HTTP 服务启动
- 默认中间件已内建：安全响应头、请求体大小限制、固定窗口限流、签名重放保护
- 已提供高层 facade：
  - `pkg/cache` JSON 读写辅助
  - `pkg/queue` 任务 facade
  - `pkg/upload` 上传 facade / uploader
  - `pkg/lock` 本地锁 / Redis 锁
  - `pkg/database` 运行时性能选项
- 健康检查路由：
  - `GET /livez`
  - `GET /readyz`

## 快速开始

```bash
git clone https://github.com/sy121718/go-mvc.git
cd go-mvc
go mod download
go run cmd/main.go
```

## 目录结构

```text
go-mvc/
├── cmd/                  # 启动入口
├── config/               # 配置解析与组件生命周期编排
├── internal/             # 业务模块、路由、中间件、任务
├── pkg/                  # 通用基础组件
├── public/               # 日志、存储、测试资源等公共目录
├── README.md
├── CLAUDE.md
└── config.yaml
```

## 生命周期

当前推荐启动顺序：

```go
if err := config.Init("config.yaml"); err != nil {
	return err
}

if err := config.InitComponents(); err != nil {
	return err
}
```

关闭时统一走：

```go
if err := config.CloseComponents(); err != nil {
	return err
}
```

## 响应约定

项目现在使用直接的 HTTP 数字状态码，不再在 `pkg` 和系统包里走统一错误码中转。

标准响应结构：

```go
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
```

使用方式：

```go
response.Success(c, data)
response.SuccessWithMessage(c, "保存成功", data)
response.ErrorWithMessage(c, 400, "请求参数错误")
response.ErrorWithMessage(c, 401, "未登录或登录已过期")
response.ErrorWithMessage(c, 404, "请求的资源不存在")
```

说明：

- `pkg/response` 只负责输出响应结构
- `pkg/response` 不维护国际化翻译
- `pkg/response` 不维护字符串错误码到状态码的映射
- 调用点直接写最终中文提示和具体数字状态码

## i18n 约定

`pkg/i18n` 仍然保留，用于直接读取数据库中的多语言文本、字典文案和 UI 文案。

但当前约定是：

- `pkg` 和系统包默认错误提示，不依赖 `i18n`
- 默认错误响应，不依赖 `i18n`
- 只有业务侧明确需要读取文案时，才直接调用 `pkg/i18n`

示例：

```go
text := i18n.GetText("ui_button_submit", "zh-CN")
httpCode := i18n.GetHttpCode("ErrAdminNotFound")
```

## 错误处理约定

- 参数校验、状态校验类错误：直接返回简单中文提示
- 底层系统错误：优先直接返回原始 `err`
- 不在 `pkg` 里做过度包装的错误翻译

不建议：

```go
return fmt.Errorf("创建日志目录失败: %w", err)
```

更倾向：

```go
return err
```

## 更多说明

- 项目级规则见 [CLAUDE.md](./CLAUDE.md)
- `pkg` 组件规则见 [pkg/CLAUDE.md](./pkg/CLAUDE.md)
- 业务模块规则见 [internal/module/CLAUDE.md](./internal/module/CLAUDE.md)
