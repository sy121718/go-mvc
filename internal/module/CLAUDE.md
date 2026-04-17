# Module 开发规范

本目录包含所有业务模块，每个模块遵循统一的目录结构和开发规范。

## 目录结构

```text
module/
└── backend/                 # 当前已启用的后台业务模块
    ├── admin/
    ├── user/
    └── ...
```

## 模块结构规范

每个模块至少包含以下核心目录：

```text
module_name/
├── router/                  # 路由目录
│   └── admin_router.go
├── handle/                  # 控制器层
│   └── admin_handle.go
├── service/                 # 业务逻辑层
│   └── admin_service.go
├── model/                   # 数据模型层
│   └── admin_model.go
├── dto/                     # 请求/响应结构（按需启用，不用时不创建空目录）
├── enums/                   # 模块级枚举和常量（按需启用，不用时不创建空目录）
└── client/                  # 对外接口（按需启用，不用时不创建空目录）
```

### 命名规范

- 外部整合层目录可以使用复数，例如 `internal/routers`
- 模块内部的小目录使用单数，例如 `router`、`handle`、`service`、`model`
- 模块内文件统一使用“模块名 + 分层名”命名：
  - `admin_router.go`
  - `admin_handle.go`
  - `admin_service.go`
  - `admin_model.go`
- 模块内包名也统一使用“模块名 + 分层名”命名，避免跨模块冲突：
  - `package adminrouter`
  - `package adminhandle`
  - `package adminservice`
  - `package adminmodel`
- 扩展层沿用同一规则：
  - `admin_dto.go` -> `package admindto`
  - `admin_enum.go` / `admin_error.go` -> `package adminenums`
  - `admin_client.go` -> `package adminclient`
  - `admin_helper.go` -> `package adminhelper`
- 不保留 `sys_` 前缀
- 不使用 `router.go`、`user.go` 这类缺少模块语义的文件名

## 多语言配置中心使用

### 1. 配置分类（category）

数据库 `sys_i18n` 表支持以下分类：

- **error** - 错误提示
- **msg** - 操作消息（成功、确认等）
- **ui** - 界面文本（按钮、标签、标题等）
- **dict** - 字典数据（状态、类型等）
- **validation** - 验证提示
- **email** - 邮件模板
- **sms** - 短信模板

### 2. 常量定义位置

**系统级常量**（所有模块通用）：

```text
internal/common/enums/
├── errors.go      # 错误码（ErrSystemError, ErrInvalidParams...）
├── messages.go    # 操作消息（MsgSaveSuccess, MsgDeleteSuccess...）
├── ui.go          # 界面文本（按钮、标签、标题等）
└── dict.go        # 字典数据（状态、类型等）
```

**模块级常量**（模块私有，按需启用）：

```text
internal/module/backend/admin/enums/
├── admin_error.go
├── admin_msg.go
└── admin_dict.go
```

### 3. 在 Handle 层使用

```go
package adminhandle

import (
    "go-mvc/internal/common/enums"
    adminEnums "go-mvc/internal/module/backend/admin/enums"
    "go-mvc/pkg/i18n"
    "go-mvc/pkg/response"

    "github.com/gin-gonic/gin"
)

type AdminHandle struct {
    service *service.AdminService
}

func (h *AdminHandle) Create(c *gin.Context) {
    var req dto.CreateAdminReq
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, enums.ErrInvalidBody)
        return
    }

    if err := h.service.Create(&req); err != nil {
        response.Error(c, err.Error())
        return
    }

    response.SuccessWithMessage(c, enums.MsgSaveSuccess, nil)
}

func (h *AdminHandle) GetUIText(c *gin.Context) {
    lang := c.GetHeader("Accept-Language")
    if lang == "" {
        lang = "zh-CN"
    }

    buttonText := i18n.Get("ui_button_submit", lang)
    titleText := i18n.Get("ui_admin_title", lang)

    response.Success(c, gin.H{
        "button": buttonText,
        "title":  titleText,
    })
}
```

### 4. 在 Service 层使用

Service 层返回错误码，不返回具体错误信息：

```go
package adminservice

import (
    "errors"
    "go-mvc/internal/common/enums"
    adminEnums "go-mvc/internal/module/backend/admin/enums"
)

type AdminService struct {
    model *model.AdminModel
}

func (s *AdminService) Create(req *dto.CreateAdminReq) error {
    exists, err := s.model.ExistsByUsername(req.Username)
    if err != nil {
        return errors.New(enums.ErrDBQueryError)
    }
    if exists {
        return errors.New(adminEnums.ErrAdminExists)
    }

    if err := s.model.Insert(req); err != nil {
        return errors.New(enums.ErrDBQueryError)
    }

    return nil
}
```

### 5. 响应方法

```go
response.Success(c, data)
response.Success(c)
response.SuccessWithMessage(c, enums.MsgSaveSuccess, data)
response.Error(c, enums.ErrSystemError)
response.Error(c, adminEnums.ErrAdminNotFound)
response.ParamError(c)
```

### 6. 语言切换

客户端可通过以下方式指定语言：

- 方式1：HTTP Header（推荐）

```text
Accept-Language: en-US
Accept-Language: zh-CN
```

- 方式2：Query 参数

```text
GET /api/admin/list?lang=en-US
GET /api/admin/list?lang=zh-CN
```

优先级：`Accept-Language Header` > `Query 参数` > `默认 zh-CN`

### 7. 直接获取多语言文本

```go
import "go-mvc/pkg/i18n"

text := i18n.Get("ui_button_submit", "zh-CN")
httpCode := i18n.GetHttpCode("ErrAdminNotFound")
```

## 开发规范

### 1. 路由定义

- 只使用 GET 和 POST 两种请求方法
- GET 用于查询操作
- POST 用于数据变更操作

```go
package adminrouter

import "github.com/gin-gonic/gin"

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

### 2. Handle 层

- 负责请求参数验证和响应
- 调用 Service 层处理业务逻辑
- 使用 response 包统一响应格式
- Service 层返回的 error 直接传给 response.Error

### 3. Service 层

- 处理业务逻辑
- 调用 Model 层操作数据
- 返回错误码（使用 errors.New(错误码常量)）
- 可以调用其他模块的 Client

### 4. Model 层

- 负责数据库操作
- 定义数据结构
- 不包含业务逻辑

### 5. DTO 层

- 定义请求和响应结构
- 使用 Gin binding 标签做验证

### 6. Helper 层

- Service 层的辅助工具按功能拆分到 `service/helper/`
- 只有存在实际复用时再创建 helper，不做空目录占位
