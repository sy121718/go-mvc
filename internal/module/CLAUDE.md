# internal/module 开发规范

`internal/module/` 存放业务模块代码。模块开发要遵循现有目录结构，并与当前响应/错误处理方式保持一致。

## 目录结构

```text
module_name/
├── router/
├── handle/
├── service/
├── model/
├── dto/      # 按需启用
├── enums/    # 按需启用
└── client/   # 按需启用
```

## 命名规则

- 文件名统一使用“模块名 + 分层名”
- 例如：
  - `admin_router.go`
  - `admin_handle.go`
  - `admin_service.go`
  - `admin_model.go`

## 分层职责

### router

- 注册模块路由
- 只负责路由组织

### handle

- 参数绑定
- 基本参数校验
- 调用 service
- 输出统一响应

### service

- 处理业务逻辑
- 组合 model / client / helper
- 返回 `error`

### model

- 只负责数据库访问
- 不放业务逻辑

## 响应约定

当前项目不再推荐在 handler 里走统一错误码中转。

推荐：

```go
response.Success(c, data)
response.SuccessWithMessage(c, "保存成功", data)
response.ErrorWithMessage(c, 400, "请求参数错误")
response.ErrorWithMessage(c, 401, "未登录或登录已过期")
response.ErrorWithMessage(c, 404, "数据不存在")
```

不再推荐：

```go
response.Error(c, enums.ErrInvalidBody)
response.SuccessWithMessage(c, enums.MsgSaveSuccess, data)
```

## handler 示例

```go
func (h *AdminHandle) Create(c *gin.Context) {
	var req dto.CreateAdminReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 400, "请求体格式错误")
		return
	}

	if err := h.service.Create(&req); err != nil {
		response.ErrorWithMessage(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "保存成功", nil)
}
```

## service 约定

service 直接返回错误，不要求统一返回错误码字符串。

推荐：

```go
func (s *AdminService) Create(req *dto.CreateAdminReq) error {
	exists, err := s.model.ExistsByUsername(req.Username)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("用户名已存在")
	}

	if err := s.model.Insert(req); err != nil {
		return err
	}

	return nil
}
```

## i18n 约定

业务模块如果明确需要读取字典、UI 文案或业务文案，可以直接调用 `pkg/i18n`。

例如：

```go
text := i18n.GetText("ui_button_submit", "zh-CN")
httpCode := i18n.GetHttpCode("ErrAdminNotFound")
```

但默认响应不依赖 i18n：

- 默认错误返回直接写中文提示
- 默认状态码直接写数字码

## 路由约定

- 只用 `GET` 和 `POST`
- `GET` 用于查询
- `POST` 用于新增、修改、删除、状态变化

示例：

```go
func SetupAdminRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	{
		admin.GET("/list", handle.List)
		admin.GET("/detail", handle.Detail)
		admin.POST("/create", handle.Create)
		admin.POST("/update", handle.Update)
		admin.POST("/delete", handle.Delete)
	}
}
```

## 错误处理原则

- 参数错误：直接返回明确中文提示
- 业务规则错误：直接返回明确中文提示
- 底层系统错误：优先直接返回原始 `err`
- 不要为了统一格式再多包一层复杂文案
