# Module 开发规范

本目录包含所有业务模块，每个模块遵循统一的目录结构和开发规范。

## 目录结构

```
module/
├── backend/           # 后台管理模块
│   ├── admin/        # 管理员模块
│   ├── user/         # 用户管理模块
│   └── ...
├── frontend/         # 前台业务模块
│   ├── user/        # 用户模块
│   ├── order/       # 订单模块
│   └── ...
└── common/          # 公共模块
    ├── enums/       # 系统级枚举和常量
    └── ...
```

## 模块结构规范

每个模块必须包含以下目录：

```
module_name/
├── router.go         # 路由定义
├── handle/          # 控制器层
├── service/         # 业务逻辑层
│   └── helper/      # 业务辅助工具
├── model/           # 数据模型层
├── dto/             # 数据传输对象
├── enums/           # 模块级枚举和常量
└── client/          # 对外接口（微服务迁移用）
```

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
```
internal/common/enums/
├── errors.go      # 错误码（ErrSystemError, ErrInvalidParams...）
├── messages.go    # 操作消息（MsgSaveSuccess, MsgDeleteSuccess...）
├── ui.go          # 界面文本（按钮、标签、标题等）
└── dict.go        # 字典数据（状态、类型等）
```

**模块级常量**（模块私有）：
```
internal/module/backend/admin/enums/
├── admin_error.go    # 管理员模块错误码
├── admin_msg.go      # 管理员模块消息
└── admin_dict.go     # 管理员模块字典
```

### 3. 在 Handle 层使用

```go
package handle

import (
    "go-mvc/internal/common/enums"
    adminEnums "go-mvc/internal/module/backend/admin/enums"
    "go-mvc/pkg/response"
    "go-mvc/pkg/i18n"
    "github.com/gin-gonic/gin"
)

func (h *AdminHandle) Create(c *gin.Context) {
    var req dto.CreateAdminReq
    if err := c.ShouldBindJSON(&req); err != nil {
        // 使用系统级错误码
        response.Error(c, enums.ErrInvalidBody)
        return
    }

    // Service 层返回错误码
    if err := h.service.Create(&req); err != nil {
        // err 是错误码字符串，直接传给 response
        response.Error(c, err.Error())
        return
    }

    // 使用系统级成功消息
    response.SuccessWithMessage(c, enums.MsgSaveSuccess, nil)
}

// 获取界面文本
func (h *AdminHandle) GetUIText(c *gin.Context) {
    lang := c.GetHeader("Accept-Language")
    if lang == "" {
        lang = "zh-CN"
    }
    
    // 直接获取多语言文本
    buttonText := i18n.Get("ui_button_submit", lang)
    titleText := i18n.Get("ui_admin_title", lang)
    
    response.Success(c, gin.H{
        "button": buttonText,
        "title": titleText,
    })
}
```

### 4. 在 Service 层使用

Service 层返回错误码，不返回具体错误信息：

```go
package service

import (
    "errors"
    "go-mvc/internal/common/enums"
    adminEnums "go-mvc/internal/module/backend/admin/enums"
)

type AdminService struct {
    model *model.AdminModel
}

func (s *AdminService) Create(req *dto.CreateAdminReq) error {
    // 业务验证 - 检查用户名是否存在
    exists, err := s.model.ExistsByUsername(req.Username)
    if err != nil {
        // 数据库错误，返回系统级错误码
        return errors.New(enums.ErrDBQueryError)
    }
    if exists {
        // 业务错误，返回模块级错误码
        return errors.New(adminEnums.ErrAdminExists)
    }

    // 执行创建
    if err := s.model.Insert(req); err != nil {
        return errors.New(enums.ErrDBQueryError)
    }

    return nil
}

func (s *AdminService) GetByID(id int) (*dto.AdminResp, error) {
    admin, err := s.model.FindByID(id)
    if err != nil {
        return nil, errors.New(enums.ErrDBQueryError)
    }
    if admin == nil {
        // 数据不存在，返回模块级错误码
        return nil, errors.New(adminEnums.ErrAdminNotFound)
    }

    return &dto.AdminResp{
        ID:       admin.ID,
        Username: admin.Username,
    }, nil
}

// 复杂业务逻辑示例
func (s *AdminService) UpdateStatus(id int, status int) error {
    // 先检查是否存在
    admin, err := s.model.FindByID(id)
    if err != nil {
        return errors.New(enums.ErrDBQueryError)
    }
    if admin == nil {
        return errors.New(adminEnums.ErrAdminNotFound)
    }

    // 业务规则验证
    if admin.IsSuper && status == 0 {
        // 超级管理员不能禁用
        return errors.New(adminEnums.ErrCannotDisableSuper)
    }

    // 执行更新
    if err := s.model.UpdateStatus(id, status); err != nil {
        return errors.New(enums.ErrDBQueryError)
    }

    return nil
}
```

### 5. 响应方法

```go
// 成功响应（带数据）
response.Success(c, data)

// 成功响应（无数据）
response.Success(c)

// 成功响应（自定义消息码）
response.SuccessWithMessage(c, enums.MsgSaveSuccess, data)

// 错误响应（自动获取多语言消息和HTTP状态码）
response.Error(c, enums.ErrSystemError)
response.Error(c, adminEnums.ErrAdminNotFound)

// 参数错误（快捷方法）
response.ParamError(c)
```

### 6. 语言切换

客户端可通过以下方式指定语言：

**方式1：HTTP Header（推荐）**
```
Accept-Language: en-US
Accept-Language: zh-CN
```

**方式2：Query 参数**
```
GET /api/admin/list?lang=en-US
GET /api/admin/list?lang=zh-CN
```

优先级：`Accept-Language Header` > `Query 参数` > `默认 zh-CN`

### 7. 数据库配置

在 `sys_i18n` 表中添加新的配置：

```sql
-- 错误码
INSERT INTO sys_i18n (item_key, lang, item_value, http_code, category, status) VALUES
('ErrAdminNotFound', 'zh-CN', '管理员不存在', 404, 'error', 1),
('ErrAdminNotFound', 'en-US', 'Admin not found', 404, 'error', 1);

-- 操作消息
INSERT INTO sys_i18n (item_key, lang, item_value, http_code, category, status) VALUES
('MsgAdminCreated', 'zh-CN', '管理员创建成功', 200, 'msg', 1),
('MsgAdminCreated', 'en-US', 'Admin created successfully', 200, 'msg', 1);

-- 界面文本
INSERT INTO sys_i18n (item_key, lang, item_value, http_code, category, status) VALUES
('ui_button_submit', 'zh-CN', '提交', 200, 'ui', 1),
('ui_button_submit', 'en-US', 'Submit', 200, 'ui', 1);

-- 字典数据
INSERT INTO sys_i18n (item_key, lang, item_value, http_code, category, status) VALUES
('dict_status_active', 'zh-CN', '启用', 200, 'dict', 1),
('dict_status_active', 'en-US', 'Active', 200, 'dict', 1);
```

系统会自动在 10 秒内刷新缓存，无需重启。

### 8. 直接获取多语言文本

除了通过 response 包自动处理，也可以直接获取多语言文本：

```go
import "go-mvc/pkg/i18n"

// 获取指定语言的文本
text := i18n.Get("ui_button_submit", "zh-CN")

// 获取HTTP状态码
httpCode := i18n.GetHttpCode("ErrAdminNotFound")

// 手动刷新缓存
if err := i18n.Reload(); err != nil {
    log.Printf("重新加载多语言配置失败: %v", err)
}
```

## 开发规范

### 1. 路由定义 (router.go)

- 只使用 GET 和 POST 两种请求方法
- GET 用于查询操作
- POST 用于数据变更操作

```go
package admin

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup) {
    admin := r.Group("/admin")
    {
        // GET - 查询类
        admin.GET("/list", handle.List)
        admin.GET("/detail", handle.GetDetail)
        
        // POST - 变更类
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

```go
package handle

import (
    "go-mvc/internal/common/enums"
    adminEnums "go-mvc/internal/module/backend/admin/enums"
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
    
    // Service 返回错误码，直接传给 response
    if err := h.service.Create(&req); err != nil {
        response.Error(c, err.Error())
        return
    }
    
    response.SuccessWithMessage(c, enums.MsgSaveSuccess, nil)
}
```

### 3. Service 层

- 处理业务逻辑
- 调用 Model 层操作数据
- 返回错误码（使用 errors.New(错误码常量)）
- 可以调用其他模块的 Client

```go
package service

import (
    "errors"
    "go-mvc/internal/common/enums"
    adminEnums "go-mvc/internal/module/backend/admin/enums"
)

type AdminService struct {
    model *model.AdminModel
}

func (s *AdminService) Create(req *dto.CreateAdminReq) error {
    // 业务验证
    exists, err := s.model.ExistsByUsername(req.Username)
    if err != nil {
        return errors.New(enums.ErrDBQueryError)
    }
    if exists {
        return errors.New(adminEnums.ErrAdminExists)
    }

    // 执行创建
    if err := s.model.Insert(req); err != nil {
        return errors.New(enums.ErrDBQueryError)
    }

    return nil
}
```

### 4. Model 层

- 负责数据库操作
- 定义数据结构
- 不包含业务逻辑

```go
package model

type Admin struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    // ...
}

type AdminModel struct{}

func (m *AdminModel) Insert(admin *Admin) error {
    // 数据库操作
    return nil
}
```

### 5. DTO 层

- 定义请求和响应结构
- 使用 Gin binding 标签做验证

```go
package dto

type CreateAdminReq struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
}

type AdminResp struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
}
```

### 6. Helper 层

Service 层的辅助工具，按功能分类：

- **formatter** - 数据格式化
- **validator** - 业务验证
- **analyzer** - 数据分析
- **matcher** - 数据匹配

```go
package helper

// formatter.go
func FormatPhone(phone string) string {
    // 格式化手机号
}

// validator.go
func ValidateIDCard(idCard string) bool {
    // 验证身份证号
}
```

## 模块间调用

使用 Client 层实现模块间调用：

```go
// internal/client/user_client.go
package client

type UserClient struct {
    service *user.UserService
}

func (c *UserClient) GetUserByID(id int) (*dto.UserResp, error) {
    return c.service.GetByID(id)
}

// 在其他模块中使用
import "go-mvc/internal/client"

userInfo, err := client.User.GetUserByID(userID)
```

## 注意事项

1. 每个模块完全独立，不直接依赖其他模块
2. 跨模块调用必须通过 Client 层
3. 模块内的 Model 是私有的，不对外暴露
4. 所有响应使用 response 包统一格式
5. 所有文本（错误、消息、界面、字典）使用多语言配置中心
6. 常量定义：系统级放 internal/common/enums，模块级放模块的 enums 目录
