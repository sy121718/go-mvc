package admindto

import adminmodel "go-mvc/internal/module/backend/admin/model"

type ListReq struct {
	Page      *int   `form:"page" json:"page" binding:"omitempty,gte=1" validate:"omitempty,gte=1"`
	Limit     *int   `form:"page_size" json:"page_size" binding:"omitempty,gte=1,lte=100" validate:"omitempty,gte=1,lte=100"`
	Email     string `form:"email" json:"email" binding:"omitempty,email_strict,max=100" validate:"omitempty,email_strict,max=100"`
	Name      string `form:"name" json:"name" binding:"omitempty,max=50" validate:"omitempty,max=50"`
	Phone     string `form:"phone" json:"phone" binding:"omitempty,max=20" validate:"omitempty,max=20"`
	Status    *int   `form:"status" json:"status"`
	SortField string `form:"sort_field" json:"sort_field" binding:"omitempty,oneof=id name email phone status create_time" validate:"omitempty,oneof=id name email phone status create_time"`
	SortOrder string `form:"sort_order" json:"sort_order" binding:"omitempty,oneof=asc desc" validate:"omitempty,oneof=asc desc"`
}

type ListResp struct {
	Total int64                    `json:"total"`
	List  []adminmodel.AdminEntity `json:"list"`
}
type SaveReq struct {
	ID       *int   `json:"id" binding:"omitempty,gte=1" validate:"omitempty,gte=1"`
	Email    string `json:"email" binding:"required,email_strict,max=100" validate:"required,email_strict,max=100"`
	Name     string `json:"name" binding:"required,max=50" validate:"required,max=50"`
	Phone    string `json:"phone" binding:"omitempty,max=20" validate:"omitempty,max=20"`
	Status   int    `json:"status" binding:"required,oneof=0 1" validate:"required,oneof=0 1"`
	Password string `json:"password" binding:"omitempty,min=6,max=100" validate:"omitempty,min=6,max=100"`
}

type SaveResp struct {
	ID int `json:"id"`
}
