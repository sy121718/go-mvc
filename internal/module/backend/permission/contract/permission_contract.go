package permissioncontract

import (
	"context"

	permissiondto "go-mvc/internal/module/backend/permission/dto"
)

// PermissionService 定义权限点模块对外暴露的业务能力。
type PermissionService interface {
	// List 权限点列表（按模块筛选）。
	List(ctx context.Context, req *permissiondto.ListReq) (*permissiondto.ListResp, error)
}
