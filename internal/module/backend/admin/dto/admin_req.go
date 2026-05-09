// Package admindto 管理员模块数据传输对象
//
// 该包定义了管理员模块在接口层与业务层之间传递的请求和响应数据结构。
// 每个结构体对应一个业务操作的输入输出格式，包含参数校验规则（binding/validate tag），
// 确保进入 service 层的数据是合法、完整的。
package admindto

// ListReq 管理员列表查询请求参数
//
// 支持分页查询和多条件模糊搜索，所有查询条件均为可选（omitempty），
// 当条件为空时不做过滤，返回全部数据。
//
// 字段说明：
//   - Page / Limit：控制分页，Limit 上限 100 条，防止一次拉取过多数据
//   - Email / Name / Phone：各自按 LIKE 模糊匹配，分别有长度限制防 SQL 注入
//   - Status：按精确值过滤，nil 表示不过滤状态
//   - SortField / SortOrder：排序控制，SortField 限定可排序字段白名单
type ListReq struct {
	Page      *int   `form:"page" json:"page" binding:"omitempty,gte=1" validate:"omitempty,gte=1"`                                                                                           // 页码，从 1 开始，不传则默认为第 1 页
	Limit     *int   `form:"page_size" json:"page_size" binding:"omitempty,gte=1,lte=100" validate:"omitempty,gte=1,lte=100"`                                                                 // 每页条数，1~100 之间，不传则默认 10 条
	Email     string `form:"email" json:"email" binding:"omitempty,email_strict,max=100" validate:"omitempty,email_strict,max=100"`                                                           // 邮箱模糊搜索，需符合邮箱格式，最长 100 字符
	Name      string `form:"name" json:"name" binding:"omitempty,max=50" validate:"omitempty,max=50"`                                                                                         // 姓名模糊搜索，最长 50 字符
	Phone     string `form:"phone" json:"phone" binding:"omitempty,max=20" validate:"omitempty,max=20"`                                                                                       // 手机号模糊搜索，最长 20 字符
	Status    *int   `form:"status" json:"status"`                                                                                                                                            // 状态精确过滤，传 1=启用 2=禁用 3=密码错误封禁；nil 表示不过滤
	SortField string `form:"sort_field" json:"sort_field" binding:"omitempty,oneof=id name email phone status create_time" validate:"omitempty,oneof=id name email phone status create_time"` // 排序字段，仅允许按 id/name/email/phone/status/create_time 排序
	SortOrder string `form:"sort_order" json:"sort_order" binding:"omitempty,oneof=asc desc" validate:"omitempty,oneof=asc desc"`                                                             // 排序方向，asc 升序 / desc 降序，不传则默认 id DESC
}

// SaveReq 管理员新增/编辑请求参数
//
// 同时用于新增和编辑两个场景：
//   - ID 为空时表示新增，创建新管理员
//   - ID 不为空时表示编辑，更新已有管理员的信息
//
// 字段说明：
//   - ID：编辑时必传，新增时不传
//   - Email / Name：必填字段
//   - Phone：可选
//   - Status：必填，0=启用 1=禁用
//   - Password：编辑时不传则保持原密码不变
type SaveReq struct {
	ID       *int   `json:"id" binding:"omitempty,gte=1" validate:"omitempty,gte=1"`                                // 管理员ID，新增时为空，编辑时必传
	Email    string `json:"email" binding:"required,email_strict,max=100" validate:"required,email_strict,max=100"` // 邮箱，必填，需合法邮箱格式，最长 100 字符
	Name     string `json:"name" binding:"required,max=50" validate:"required,max=50"`                              // 姓名，必填，最长 50 字符
	Phone    string `json:"phone" binding:"omitempty,max=20" validate:"omitempty,max=20"`                           // 手机号，可选，最长 20 字符
	Status   int    `json:"status" binding:"required,oneof=0 1" validate:"required,oneof=0 1"`                      // 状态，必填，0=启用 1=禁用
	Password string `json:"password" binding:"omitempty,min=6,max=100" validate:"omitempty,min=6,max=100"`          // 密码，新增时必传，编辑时不传则保持原密码
}