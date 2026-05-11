// Package admindto 管理员模块数据传输对象
//
// 该包定义了管理员模块在接口层与业务层之间传递的请求和响应数据结构。
// 每个结构体对应一个业务操作的输入输出格式，包含参数校验规则（binding/validate tag），
// 确保进入 service 层的数据是合法、完整的。
package admindto

import adminmodel "go-mvc/internal/module/backend/admin/model"

// ListResp 管理员列表查询响应
//
// 返回符合查询条件的总记录数和当前页的数据列表。
type ListResp struct {
	Total int64                    `json:"total"` // 符合条件的总记录数，用于前端分页组件
	List  []adminmodel.AdminEntity `json:"list"`  // 当前页的管理员数据列表
}

// CreateResp 管理员新增响应
type CreateResp struct {
	ID       uint64 `json:"id"` // 新增的管理员 ID，表示添加成功
	Username string `json:"username"`
}

// LoginResp 管理员登录响应
//
// 登录成功后返回 token 和用户基本信息，前端会根据此信息写入 localStorage。
type LoginResp struct {
	AccessToken  string   `json:"accessToken"`                                           // JWT token，后续请求带在 Authorization 头
	RefreshToken string   `json:"refreshToken"`                                          // 刷新 token，accessToken 过期后用它续期
	Expires      string   `json:"expires"`                                               // token 过期时间，格式 "2006/01/02 15:04:05"
	Username     string   `json:"username"`                                              // 登录账号
	Nickname     string   `json:"nickname"`                                              // 显示昵称（对应 admin_entity.name）
	Avatar       string   `json:"avatar"`                                                // 头像URL
	Roles        []string `json:"roles"`                                                 // 角色列表，用于前端页面级权限判断
	Permissions  []string `json:"permissions"`                                           // 按钮权限列表，前端按钮级别显隐控制
}
