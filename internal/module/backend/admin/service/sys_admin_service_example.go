package service

import (
	"errors"
	"go-mvc/internal/module/backend/admin/model"
	"go-mvc/pkg/auth"
	"go-mvc/pkg/database"
	"time"

	"gorm.io/gorm"
)

/*
AdminService - 业务逻辑层
===========================================
职责：
1. 业务流程编排
2. 业务规则验证
3. 直接调用 GORM 操作数据
4. 调用 Model 的自身规则方法
5. 调用其他 Service 或外部服务

不负责：
❌ HTTP 请求处理
❌ 参数绑定
❌ 响应格式化
*/

// AdminService 管理员服务
type AdminService struct{}

// NewAdminService 创建服务实例
func NewAdminService() *AdminService {
	return &AdminService{}
}

// ====== 登录相关 ======

// Login 登录
func (s *AdminService) Login(username, password string) (*model.SysAdmin, string, error) {
	// 1. 查询用户 - 直接调用 GORM
	var admin model.SysAdmin
	err := database.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("用户不存在")
		}
		return nil, "", errors.New("查询失败")
	}

	// 2. 检查状态 - 调用 Model 的自身规则
	if !admin.IsActive() {
		return nil, "", errors.New("账号已被禁用")
	}

	// 3. 验证密码 - 调用 Model 的自身规则
	if !admin.CheckPassword(password) {
		return nil, "", errors.New("密码错误")
	}

	// 4. 生成 Token
	token, err := auth.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		return nil, "", errors.New("生成Token失败")
	}

	// 5. 更新登录时间 - 直接调用 GORM
	database.DB.Model(&admin).Update("updated_at", time.Now())

	return &admin, token, nil
}

// Logout 登出
func (s *AdminService) Logout(token string) error {
	// TODO: 将 Token 加入黑名单
	return nil
}

// ====== CRUD 操作 ======

// Create 创建管理员
func (s *AdminService) Create(username, password, email string) (*model.SysAdmin, error) {
	// 1. 检查用户是否存在 - 直接调用 GORM
	var existAdmin model.SysAdmin
	err := database.DB.Where("username = ?", username).First(&existAdmin).Error
	if err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 2. 创建实例
	admin := &model.SysAdmin{
		Username: username,
		Email:    email,
		Status:   1,
	}

	// 3. 设置密码 - 调用 Model 的自身规则
	if err := admin.SetPassword(password); err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 4. 验证邮箱 - 调用 Model 的自身规则
	if !admin.ValidateEmail() {
		return nil, errors.New("邮箱格式错误")
	}

	// 5. 保存到数据库 - 直接调用 GORM
	if err := database.DB.Create(admin).Error; err != nil {
		return nil, errors.New("创建失败")
	}

	return admin, nil
}

// GetByID 根据ID查询
func (s *AdminService) GetByID(id int64) (*model.SysAdmin, error) {
	var admin model.SysAdmin
	err := database.DB.First(&admin, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("管理员不存在")
		}
		return nil, errors.New("查询失败")
	}
	return &admin, nil
}

// GetByUsername 根据用户名查询
func (s *AdminService) GetByUsername(username string) (*model.SysAdmin, error) {
	var admin model.SysAdmin
	err := database.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("管理员不存在")
		}
		return nil, errors.New("查询失败")
	}
	return &admin, nil
}

// List 获取列表
func (s *AdminService) List(page, pageSize int, keyword string) ([]model.SysAdmin, int64, error) {
	var admins []model.SysAdmin
	var total int64

	// 构建查询
	query := database.DB.Model(&model.SysAdmin{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&admins).Error
	if err != nil {
		return nil, 0, errors.New("查询失败")
	}

	return admins, total, nil
}

// Update 更新
func (s *AdminService) Update(id int64, updates map[string]interface{}) error {
	// 检查是否存在
	var admin model.SysAdmin
	if err := database.DB.First(&admin, id).Error; err != nil {
		return errors.New("管理员不存在")
	}

	// 更新 - 直接调用 GORM
	result := database.DB.Model(&admin).Updates(updates)
	if result.Error != nil {
		return errors.New("更新失败")
	}

	return nil
}

// UpdatePassword 更新密码
func (s *AdminService) UpdatePassword(id int64, oldPassword, newPassword string) error {
	// 1. 查询管理员
	var admin model.SysAdmin
	if err := database.DB.First(&admin, id).Error; err != nil {
		return errors.New("管理员不存在")
	}

	// 2. 验证旧密码 - 调用 Model 的自身规则
	if !admin.CheckPassword(oldPassword) {
		return errors.New("旧密码错误")
	}

	// 3. 设置新密码 - 调用 Model 的自身规则
	if err := admin.SetPassword(newPassword); err != nil {
		return errors.New("密码加密失败")
	}

	// 4. 更新 - 直接调用 GORM
	if err := database.DB.Save(&admin).Error; err != nil {
		return errors.New("更新失败")
	}

	return nil
}

// Delete 删除
func (s *AdminService) Delete(id int64) error {
	// 直接调用 GORM 删除
	result := database.DB.Delete(&model.SysAdmin{}, id)
	if result.Error != nil {
		return errors.New("删除失败")
	}

	if result.RowsAffected == 0 {
		return errors.New("管理员不存在")
	}

	return nil
}

// UpdateStatus 更新状态
func (s *AdminService) UpdateStatus(id int64, status int) error {
	// 检查是否存在
	var admin model.SysAdmin
	if err := database.DB.First(&admin, id).Error; err != nil {
		return errors.New("管理员不存在")
	}

	// 更新状态 - 直接调用 GORM
	result := database.DB.Model(&admin).Update("status", status)
	if result.Error != nil {
		return errors.New("更新失败")
	}

	return nil
}

// ====== 批量操作 ======

// BatchDelete 批量删除
func (s *AdminService) BatchDelete(ids []int64) error {
	// 直接调用 GORM 批量删除
	result := database.DB.Delete(&model.SysAdmin{}, ids)
	if result.Error != nil {
		return errors.New("批量删除失败")
	}

	return nil
}

// BatchUpdateStatus 批量更新状态
func (s *AdminService) BatchUpdateStatus(ids []int64, status int) error {
	// 直接调用 GORM 批量更新
	result := database.DB.Model(&model.SysAdmin{}).Where("id IN ?", ids).Update("status", status)
	if result.Error != nil {
		return errors.New("批量更新失败")
	}

	return nil
}

// ====== 统计相关 ======

// CountByStatus 按状态统计
func (s *AdminService) CountByStatus(status int) (int64, error) {
	var count int64
	err := database.DB.Model(&model.SysAdmin{}).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return 0, errors.New("统计失败")
	}
	return count, nil
}

// GetStatistics 获取统计信息
func (s *AdminService) GetStatistics() (map[string]interface{}, error) {
	var total, active, disabled int64

	// 总数
	database.DB.Model(&model.SysAdmin{}).Count(&total)

	// 激活数
	database.DB.Model(&model.SysAdmin{}).Where("status = 1").Count(&active)

	// 禁用数
	database.DB.Model(&model.SysAdmin{}).Where("status = 0").Count(&disabled)

	return map[string]interface{}{
		"total":    total,
		"active":   active,
		"disabled": disabled,
	}, nil
}
