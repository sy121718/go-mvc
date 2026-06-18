package permissiondto

// PermissionItem 权限点列表项
type PermissionItem struct {
	ID             uint64 `json:"id"`
	PermissionCode string `json:"permission_code"`
	PermissionName string `json:"permission_name"`
	Module         string `json:"module"`
	APIPath        string `json:"api_path"`
	APIMethod      string `json:"api_method"`
	Status         int    `json:"status"`
}

// ListResp 列表查询响应
type ListResp struct {
	Total int64 `json:"total"`
}
