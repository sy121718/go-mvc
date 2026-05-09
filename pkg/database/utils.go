package database

import "gorm.io/gorm"

// IsFieldExists 检查数据库表中是否存在指定字段值的记录
//
// 参数:
//   - db: GORM 数据库连接实例
//   - model: 模型对象，用于指定要查询的表（例如 &User{} 或 User{})
//   - field: 要检查的字段名（数据库列名）
//   - value: 要检查的字段值
//
// 返回值:
//   - bool: true 表示存在匹配的记录，false 表示不存在
//   - error: 数据库查询错误，查询成功时返回 nil
//
// 注意: 该函数会自动过滤软删除的记录（deleted_time IS NULL）
//
// 使用示例:
//
//	exists, err := ExistsByField(db, &User{}, "email", "user@example.com")
//	if err != nil {
//	    // 处理错误
//	}
//	if exists {
//	    // 记录已存在
//	}
func IsFieldExists(db *gorm.DB, model interface{}, field string, value interface{}) (bool, error) {
	var count int64
	err := db.Model(model).
		Where(field+" = ? AND deleted_time IS NULL", value).
		Count(&count).Error
	return count > 0, err
}
