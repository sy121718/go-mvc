package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

const(
	AdminStatusActive = 1
	AdminStatusInactive = 0
)

// Admin 管理员模型（充血模型 - 字段 + 自动规则 + 状态判断）
type Admin struct {
	ID                uint       `gorm:"primaryKey"`                                                    // 主键
	Username          string     `gorm:"type:varchar(50);uniqueIndex" validate:"required,min=3,max=20"` // 用户名（唯一）
	Password          string     `gorm:"type:varchar(100)" validate:"required,min=6"`                   // 密码
	Name              string     `gorm:"type:varchar(50)" validate:"required"`                          // 姓名
	Avatar            string     `gorm:"type:varchar(255)"`                                             // 头像
	Email             string     `gorm:"type:varchar(100);index" validate:"email"`                      // 邮箱（索引）
	Phone             string     `gorm:"type:varchar(20);index" validate:"omitempty,phone"`             // 手机号（索引）
	Status            int        `gorm:"type:tinyint(4);default:1;index" validate:"oneof=0 1"`          // 状态：1启用 0禁用（索引）
	IsAdmin           int        `gorm:"type:tinyint(4);default:0" protected:"true"`                    // 是否超级管理员：1是 0否（受保护字段）
	LoginFailureCount int        `gorm:"type:smallint unsigned;default:0"`                              // 登录失败次数
	LockedUntilTime   *time.Time `gorm:"type:datetime(3)"`                                              // 锁定截止时间（可空）
	LastFailureTime   *time.Time `gorm:"type:datetime(3)"`                                              // 最后失败时间（可空）
	RegisterIP        string     `gorm:"type:varchar(50)"`                                              // 注册IP
	RegisterLocation  string     `gorm:"type:varchar(100)"`                                             // 注册地点
	LastLoginIP       string     `gorm:"type:varchar(50)"`                                              // 最后登录IP
	LastLoginLocation string     `gorm:"type:varchar(100)"`                                             // 最后登录地点
	LastLoginISP      string     `gorm:"type:varchar(50)"`                                              // 最后登录运营商
	LastLoginTime     *time.Time `gorm:"type:datetime(3)"`                                              // 最后登录时间（可空）
	CreateBy          uint       `gorm:"type:bigint unsigned"`                                          // 创建人ID
	CreateTime        *time.Time `gorm:"type:datetime(3)"`                                              // 创建时间
	UpdateBy          uint       `gorm:"type:bigint unsigned"`                                          // 更新人ID
	UpdateTime        *time.Time `gorm:"type:datetime(3)"`                                              // 更新时间
	DeletedTime       *time.Time `gorm:"type:datetime(3)"`                                              // 软删除时间
}

// TableName 指定表名
func (a *Admin) TableName() string {
	return "sys_admin"
}

// ==================== 状态判断规则（纯逻辑）====================

// CanLogin 是否可以登录
func (a *Admin) CanLogin() bool {
	return a.Status == 1 && !a.IsLocked()
}

// IsLocked 是否被锁定
func (a *Admin) IsLocked() bool {
	if a.LockedUntilTime == nil {
		return false
	}
	return time.Now().Before(*a.LockedUntilTime)
}

// IsActive 是否启用状态
func (a *Admin) IsActive() bool {
	return a.Status == 1
}

// IsSuperAdmin 是否超级管理员
func (a *Admin) IsSuperAdmin() bool {
	return a.IsAdmin == 1
}

// CanModifyField 是否允许修改指定字段（受保护字段不可外部修改）
func (a *Admin) CanModifyField(field string) bool {
	protectedFields := map[string]bool{
		"IsAdmin": true, // 超管字段受保护
	}
	return !protectedFields[field]
}

// ==================== GORM 钩子（自动规则，不需要外部调用）====================

// BeforeCreate 创建前自动设置时间 + 检查超管规则
func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	// 1. 自动设置时间
	now := time.Now()
	a.CreateTime = &now
	a.UpdateTime = &now

	// 2. 如果要创建超管，检查数据库是否已有超管（规则：只能有一个）
	if a.IsAdmin == 1 {
		var count int64
		tx.Model(&Admin{}).Where("is_admin = 1").Count(&count)
		if count > 0 {
			return errors.New("系统已存在超级管理员，不能重复创建")
		}
	}
	return nil
}

// BeforeUpdate 更新前自动设置时间 + 阻止受保护字段修改
func (a *Admin) BeforeUpdate(tx *gorm.DB) error {
	// 1. 自动设置时间
	now := time.Now()
	a.UpdateTime = &now
	// 2. 阻止 IsAdmin 字段被外部修改（规则：超管字段不可外部修改）
	if tx.Statement.Changed("IsAdmin") {
		// 只有内部调用（设置标记）才允许修改，外部请求直接拒绝
		if !tx.Statement.Context.Value("allow_modify_is_admin").(bool) {
			return errors.New("不允许外部修改")
		}
	}

	return nil
}
