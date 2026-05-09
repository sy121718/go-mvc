package adminservice

import (
	"context"
	"errors"
	"strings"

	admindto "go-mvc/internal/module/backend/admin/dto"
	adminmodel "go-mvc/internal/module/backend/admin/model"
	"go-mvc/pkg/database"

	"golang.org/x/crypto/bcrypt"
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

// 新增管理员
func (s *Service) Create(ctx context.Context, req *admindto.CreateReq) (*admindto.CreateResp, error) {
	// 检查邮箱是否已存在
	var existCount int64
	if err := s.am.Query(ctx).Where("email = ? AND deleted_time IS NULL", req.Email).Count(&existCount).Error; err != nil {
		return nil, err
	}
	if existCount > 0 {
		return nil, errors.New("该邮箱已被占用")
	}
	if emailExists, err := database.IsFieldExists(s.am.Query(ctx), &adminmodel.AdminEntity{}, "email", req.Email); err != nil {
		return nil, err
	} else if emailExists {
		return nil, errors.New("该邮箱已存在")
	}
	if nameExists, err := database.IsFieldExists(s.am.Query(ctx), &adminmodel.AdminEntity{}, "name", req.Username); err != nil {
		return nil, err
	} else if nameExists {
		return nil, errors.New("用户名已存在，请修改")
	}
	//可以用下面这种，就是不方便阅读，但是更好调试
	// emailExists, err := database.IsFieldExists(s.am.Query(ctx), &adminmodel.AdminEntity{}, "email", req.Email)
	// if err != nil {
	// 	fmt.Print("666")
	// 	return nil, err
	// }
	// if emailExists {
	// 	return nil, errors.New("该邮箱已存在")
	// }

	// 加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 构造实体
	//取值问题：数据库设计的可以为空，然后在模型里面设置的*开头的类型，只能取内存地址，不能直接放值进去
	entity := &adminmodel.AdminEntity{
		Username: req.Username,
		Password: string(hashed),
		Email:    &req.Email,
		Status:   adminmodel.AdminStatusActive,
		Name:     &req.Username, // 反正一定有值，直接写里面
	}
	//有接收的，都必须判断是否存在
	if req.Phone != "" {
		entity.Phone = &req.Phone
	}
	if req.Avatar != "" {
		entity.Avatar = &req.Avatar
	}

	if err := s.am.Query(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	res := &admindto.CreateResp{
		ID:       entity.ID,
		Username: entity.Username,
	}

	return res, nil
}
