// Package admindto 管理员模块数据传输对象
//
// 该包定义了管理员模块在接口层与业务层之间传递的请求和响应数据结构。
// 每个结构体对应一个业务操作的输入输出格式，包含参数校验规则（binding/validate tag），
// 确保进入 service 层的数据是合法、完整的。
package admindto

import (
	"time"
)

// ListResp 管理员列表查询响应
//
// 返回符合查询条件的总记录数和当前页的数据列表。
type ListResp struct {
	Total int64       `json:"total"` // 符合条件的总记录数，用于前端分页组件
	List  []AdminItem `json:"list"`  // 当前页的管理员数据列表
}

// AdminItem 列表项，只返回前端需要的字段，不暴露敏感/内部字段。
type AdminItem struct {
	ID         uint64     `json:"id"`
	Username   string     `json:"username"`
	Name       *string    `json:"name"`
	Avatar     *string    `json:"avatar"`
	Email      *string    `json:"email"`
	Phone      *string    `json:"phone"`
	Status     int        `json:"status"`      // 1启用 2禁用 3封禁
	CreateTime *time.Time `json:"create_time"` // 创建时间
}

// CreateResp 管理员新增响应
type CreateResp struct {
	ID       uint64 `json:"id"` // 新增的管理员 ID，表示添加成功
	Username string `json:"username"`
}

// LoginResp 管理员登录响应
//
// token 同时通过响应头和响应体下发，避免跨域场景下前端读不到自定义响应头。
type LoginResp struct {
	AccessToken string `json:"accessToken"` // JWT token
}

// DetailResp 管理员详情响应
//
// 一个管理员查看另一个管理员的详细信息。
// 不包含 password、login_failure_count 等内部安全字段，不含 update_by / update_time。
type DetailResp struct {
	ID                uint64     `json:"id"`
	Username          string     `json:"username"`
	Name              string     `json:"name"`
	Avatar            string     `json:"avatar"`
	Email             string     `json:"email"`
	Phone             string     `json:"phone"`
	Status            int        `json:"status"`   // 1启用 2禁用 3封禁
	IsAdmin           int        `json:"is_admin"` // 是否超管
	Roles             []any      `json:"roles"`    // 角色列表（由 service 层组装）
	Menus             []any      `json:"menus"`    // 菜单列表（由 service 层组装）
	RegisterIP        string     `json:"register_ip"`
	RegisterLocation  string     `json:"register_location"`
	LastLoginIP       string     `json:"last_login_ip"`
	LastLoginLocation string     `json:"last_login_location"`
	LastLoginTime     *time.Time `json:"last_login_time"`
	CreateBy          uint64     `json:"create_by"`
	CreateTime        *time.Time `json:"create_time"`
	Remark            string     `json:"remark"`
}

// ProfileResp 当前登录用户信息响应（从 Redis 会话或数据库获取）。
type ProfileResp struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"` //备注
	Menus    []any  `json:"menus"`
}

// CreateResp 管理员新增响应
type EditResp struct {
	ID uint64 `json:"id"` // 新增的管理员 ID，表示编辑成功
}
type DeleteResp struct {
	DeletedCount int64 `json:"deleted_count"`
}
