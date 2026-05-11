package adminservice

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	admindto "go-mvc/internal/module/backend/admin/dto"
	adminmodel "go-mvc/internal/module/backend/admin/model"
	"go-mvc/pkg/auth"
	"go-mvc/pkg/captcha"
	"go-mvc/pkg/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

// 新增管理员业务层方法-只新增，不管编辑更新
func (s *Service) Create(ctx context.Context, req *admindto.CreateReq) (*admindto.CreateResp, error) {
	// // 检查邮箱是否已存在
	// var existCount int64
	// if err := s.am.Query(ctx).Where("email = ? AND deleted_time IS NULL", req.Email).Count(&existCount).Error; err != nil {
	// 	return nil, err
	// }
	// if existCount > 0 {
	// 	return nil, errors.New("该邮箱已被占用")
	// }
	if emailExists, err := database.IsFieldExists(s.am.Query(ctx), &adminmodel.AdminEntity{}, "email", req.Email); err != nil {
		return nil, err
	} else if emailExists {
		return nil, errors.New("该邮箱已存在")
	}
	if nameExists, err := database.IsFieldExists(s.am.Query(ctx), &adminmodel.AdminEntity{}, "username", req.Username); err != nil {
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
	// *string 类型的字段不能直接赋 string 值，需要用 & 取地址
	entity := &adminmodel.AdminEntity{
		Username: req.Username,
		Password: string(hashed),
		Email:    &req.Email,
		Status:   adminmodel.AdminStatusActive,
		Name:     &req.Username, // Username 必填，直接赋值
	}
	// 只要有接收值并且可选字段的（omitempty）：有值才写，没值保持 nil → 数据库写 NULL，
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

// Login 管理员登录，返回 token 和用户信息。
//
// 流程：
//  1. 验证图形验证码
//  2. 查数据库，找不到用户返回模糊错误
//  3. 检查是否被锁定 / 被动禁用
//  4. bcrypt 对比密码
//  5. 失败：累加失败次数，连续 5 次后封禁 30 分钟
//  6. 成功：清空失败状态，记录登录 IP 和时间，生成 token 对
func (s *Service) Login(ctx context.Context, req *admindto.LoginReq, clientIP string) (*admindto.LoginResp, error) {
	// 1) 验证验证码
	var captchaSvc = captcha.Get()
	// Verify-验证验证码是否正确
	if !captchaSvc.Verify(req.CaptchaID, req.Captcha, true) {
		return nil, errors.New("验证码错误或已过期")
	}

	// 2) 按用户名查用户（区分大小写）
	var entity adminmodel.AdminEntity
	//go特色，if先写短函数，然后再定义条件，当前是捕获err，如果err不为空，那就找错误，
	// 第一层错误是发牛的nil和系统报错，第二层是捕获gorm的报错然后替代默认的err
	if err := s.am.Query(ctx).Where("(BINARY username = ? OR email = ?)", req.Username, req.Username).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 3) 检查是否被锁定
	if entity.IsLocked() {
		return nil, fmt.Errorf("账号已被锁定，请 %s 后重试",
			time.Until(*entity.LockedUntilTime).Round(time.Minute).String())
	}

	// 4) 检查是否被禁用
	if !entity.IsActive() {
		return nil, errors.New("账号已被禁用")
	}

	// 5) 密码校验
	if err := bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(req.Password)); err != nil {
		recordLoginFailure(ctx, s.am, &entity)
		return nil, errors.New("用户名或密码错误")
	}

	// 6) 登录成功，清空失败状态，记录登录信息
	now := time.Now()            //获取当前时间
	entity.LoginFailureCount = 0 //登录失败次数清空为0
	entity.LockedUntilTime = nil //清空封禁时间
	entity.LastFailureTime = nil
	entity.LastLoginTime = &now //设置登录时间
	if clientIP != "" {
		entity.LastLoginIP = &clientIP
	}
	if err := s.am.Query(ctx).Select(
		"login_failure_count", "locked_until_time", "last_failure_time",
		"last_login_time", "last_login_ip").
		Updates(&entity).Error; err != nil {
		return nil, err
	}

	// 7) 生成 token 对 传入
	accessToken, refreshToken, _, err := auth.GenerateTokenPair(int64(entity.ID), entity.Username, req.RememberMe)
	if err != nil {
		return nil, fmt.Errorf("生成 token 失败: %w", err)
	}
	nickname := *entity.Name
	avatar := ""
	if entity.Avatar != nil {
		avatar = *entity.Avatar
	}
	email := ""
	if entity.Email != nil {
		email = *entity.Email
	}

	return &admindto.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Username:     entity.Username,
		Nickname:     nickname,
		Email:        email,
		Avatar:       avatar,
	}, nil
}

// recordLoginFailure 记录登录失败：累加次数，连续 5 次封禁 30 分钟。
func recordLoginFailure(ctx context.Context, am *adminmodel.AdminModel, entity *adminmodel.AdminEntity) {
	now := time.Now()
	entity.LoginFailureCount++
	entity.LastFailureTime = &now

	if entity.LoginFailureCount >= 5 {
		entity.Status = adminmodel.AdminStatusPasswordError
		lockedUntil := now.Add(30 * time.Minute)
		entity.LockedUntilTime = &lockedUntil
	}

	am.Query(ctx).Select("login_failure_count", "last_failure_time", "status", "locked_until_time").
		Updates(entity)
}
