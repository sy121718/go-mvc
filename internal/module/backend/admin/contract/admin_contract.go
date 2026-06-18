package admincontract

import (
	"context"

	admindto "go-mvc/internal/module/backend/admin/dto"
)

// AdminService 定义 handle 层需要的 admin 业务能力（对外暴露契约）。
type AdminService interface {
	// List 管理员列表（分页 + 筛选）。
	List(ctx context.Context, req *admindto.ListReq) (*admindto.ListResp, error)
	// Login 管理员登录，返回 token 与基本信息；clientIP 用于登录日志/风控。
	Login(ctx context.Context, req *admindto.LoginReq, clientIP string) (*admindto.LoginResp, error)
	// Profile 获取当前登录管理员的个人信息。
	Profile(ctx context.Context, userID uint64) (*admindto.ProfileResp, error)
	// Create 新增管理员。
	Create(ctx context.Context, req *admindto.CreateReq) (*admindto.CreateResp, error)
	// Edit 修改管理员信息。
	Edit(ctx context.Context, req *admindto.EditReq) (*admindto.EditResp, error)
	// Detail 获取单个管理员详情。
	Detail(ctx context.Context, req *admindto.DetailReq) (*admindto.DetailResp, error)
	// Delete 删除管理员。
	Delete(ctx context.Context, req *admindto.DeleteReq) (*admindto.DeleteResp, error)
}
