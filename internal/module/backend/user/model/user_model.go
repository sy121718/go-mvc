package usermodel

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey"`                    // 主键自增
	Username  string    `gorm:"type:varchar(50);uniqueIndex"`  // 用户名，唯一索引
	Password  string    `gorm:"type:varchar(255)"`             // 密码（加密存储）
	Nickname  string    `gorm:"type:varchar(50)"`              // 昵称
	Email     string    `gorm:"type:varchar(100)"`             // 邮箱
	Phone     string    `gorm:"type:varchar(20)"`              // 手机号
	Avatar    string    `gorm:"type:varchar(255)"`             // 头像
	Status    int       `gorm:"type:tinyint;default:1"`        // 状态：1启用 0禁用
	DeptID    *uint     `gorm:"index"`                        // 部门ID（可空）
	Remark    string    `gorm:"type:varchar(255)"`            // 备注
	CreatedAt time.Time                                    // 创建时间
	UpdatedAt time.Time                                    // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index"`                  // 软删除
}

// TableName 指定表名
func (User) TableName() string {
	return "sys_user"
}
