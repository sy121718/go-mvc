package admincontract

import (
	"context"

	admindto "go-mvc/internal/module/backend/admin/dto"
)

// AdminService 定义 handle 层需要的 admin 业务能力。
type AdminService interface {
	List(ctx context.Context, req *admindto.ListReq) (*admindto.ListResp, error)
	Login(ctx context.Context, req *admindto.LoginReq, clientIP string) (*admindto.LoginResp, error)
	Profile(ctx context.Context, userID uint64) (*admindto.ProfileResp, error)
	Create(ctx context.Context, req *admindto.CreateReq) (*admindto.CreateResp, error)
	Edit(ctx context.Context, req *admindto.EditReq) (*admindto.EditResp, error)
	Detail(ctx context.Context, req *admindto.DetailReq) (*admindto.DetailResp, error)
	Delete(ctx context.Context, req *admindto.DeleteReq) (*admindto.DeleteResp, error)
}
