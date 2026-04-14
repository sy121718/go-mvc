package validate

import (
	"gorm.io/gorm"
)

// IsUnique 检查字段值是否唯一（不存在）
// model: 模型实例（用于获取表名）
// field: 字段名（如 "username"）
// value: 字段值
func IsUnique(db *gorm.DB, model any, field string, value any) bool {
	var count int64
	db.Model(model).Where(field+" = ?", value).Count(&count)
	return count == 0
}

// IsExists 检查字段值是否存在
// model: 模型实例
// field: 字段名
// value: 字段值
func IsExists(db *gorm.DB, model any, field string, value any) bool {
	var count int64
	db.Model(model).Where(field+" = ?", value).Count(&count)
	return count > 0
}

// IsUniqueExclude 排除当前记录后检查唯一性（用于更新场景）
// model: 模型实例
// field: 字段名
// value: 字段值
// excludeID: 排除的记录ID（当前记录）
func IsUniqueExclude(db *gorm.DB, model any, field string, value any, excludeID uint) bool {
	var count int64
	db.Model(model).Where(field+" = ?", value).Where("id != ?", excludeID).Count(&count)
	return count == 0
}

// IsUniqueExcludeField 排除指定字段后检查唯一性（更灵活）
// model: 模型实例
// field: 要检查的字段名
// value: 字段值
// excludeField: 排除的字段名（如 "id"）
// excludeValue: 排除的字段值
func IsUniqueExcludeField(db *gorm.DB, model any, field string, value any, excludeField string, excludeValue any) bool {
	var count int64
	db.Model(model).Where(field+" = ?", value).Where(excludeField+" != ?", excludeValue).Count(&count)
	return count == 0
}