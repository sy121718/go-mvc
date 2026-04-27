package adminservice

import (
	"context"
	"errors"
	"strings"

	admindto "go-mvc/internal/module/backend/admin/dto"
	adminmodel "go-mvc/internal/module/backend/admin/model"
)

// Service 定义了 Admin 模块的业务逻辑
type Service struct {
	am *adminmodel.AdminModel
}

// Deps 定义了 Service 依赖的组件
type Deps struct {
	AdminModel *adminmodel.AdminModel
}

// NewService 创建一个 Service 实例对象
func NewService(deps Deps) *Service {
	return &Service{
		am: deps.AdminModel,
	}
}

func (s *Service) List(c context.Context, req *admindto.ListReq) (res *admindto.ListResp, err error) {
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}

	var total int64
	page := 1
	if req.Page != nil {
		page = *req.Page
	}

	limit := 10
	if req.Limit != nil {
		limit = *req.Limit
	}

	query := s.am.Query(c).Where("deleted_time IS NULL")

	if email := strings.TrimSpace(req.Email); email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	if name := strings.TrimSpace(req.Name); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if phone := strings.TrimSpace(req.Phone); phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if err = query.Count(&total).Error; err != nil {
		return nil, err
	}

	orderClause := "id DESC"
	if req.SortField != "" && req.SortOrder != "" {
		orderClause = string(req.SortField) + " " + strings.ToUpper(string(req.SortOrder))
	}

	list := make([]adminmodel.AdminEntity, 0, limit)
	if err = query.
		Order(orderClause).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}

	res = &admindto.ListResp{
		Total: total,
		List:  list,
	}

	return res, nil
}

func (s *Service) Save(c context.Context, req *admindto.SaveReq) (res *admindto.SaveResp, err error) {

	return nil, nil
}
