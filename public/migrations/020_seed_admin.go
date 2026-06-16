// 默认管理员种子数据：用户名 admin，密码 admin123（bcrypt 加密），超管权限。
package migrations

func init() {
	registerSeed(Seed{
		Version:      "020",
		TableName:    "sys_admin",
		ConditionSQL: "SELECT COUNT(*) FROM sys_admin WHERE username = 'admin'",
		SQL: `INSERT INTO sys_admin (id, username, password, name, avatar, email, phone, status, is_admin, login_failure_count, locked_until_time, metadata, last_failure_time, register_ip, register_location, last_login_ip, last_login_location, last_login_isp, last_login_time, create_by, create_time, update_by, update_time, remark) VALUES
(1, 'admin', '$2a$10$Xk2cyVAGlTtfBtzTOiD5ae4FYHzLSz9lR1ggGx3r.I7CLCer.Hg4y', 'admin', 'https://avatars.githubusercontent.com/u/52823142', '1217189608@qq.com', NULL, 1, 1, 0, NULL, NULL, NULL, NULL, NULL, '::1', NULL, NULL, NOW(3), 0, NULL, 0, NULL, NULL)`,
	})
}
