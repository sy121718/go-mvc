// 系统管理员表：核心用户表，存储登录账号、密码（bcrypt）、
// 状态（启用/禁用/密码错误封禁）、登录信息（IP/地理位置/运营商/时间）、
// 登录失败锁定机制、扩展元数据等。
package migrations

func init() {
	register(Migration{
		Version:   "008",
		TableName: "sys_admin",
		SQL: `CREATE TABLE sys_admin (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID（唯一）',
  username varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录账号用户名',
  password varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '加密密码（如bcrypt）',
  name varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户姓名',
  avatar varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '头像URL',
  email varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '邮箱',
  phone varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '手机号',
  status tinyint NOT NULL DEFAULT '1' COMMENT '状态：1启用、2禁用、3密码错误封禁',
  is_admin tinyint NOT NULL DEFAULT '0' COMMENT '是否管理员：0否、1是',
  login_failure_count smallint unsigned NOT NULL DEFAULT '0' COMMENT '连续登录失败次数（达到阈值后临时锁定）',
  locked_until_time datetime(3) DEFAULT NULL COMMENT '封禁至（NULL表示未封禁）',
  metadata json DEFAULT NULL COMMENT '扩展元数据',
  last_failure_time datetime(3) DEFAULT NULL COMMENT '最后一次登录失败时间',
  register_ip varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '注册IP地址',
  register_location varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '注册地理位置（如：北京市-联通）',
  last_login_ip varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后登录IP',
  last_login_location varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后登录地理位置',
  last_login_isp varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后登录网络运营商',
  last_login_time datetime(3) DEFAULT NULL COMMENT '最后登录时间',
  create_by bigint unsigned NOT NULL COMMENT '创建人ID（0=系统创建）',
  create_time datetime(3) DEFAULT NULL COMMENT '创建时间',
  update_by bigint unsigned NOT NULL COMMENT '更新人ID（0=系统更新）',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  remark varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY uk_username (username) USING BTREE COMMENT '用户名唯一索引',
  KEY idx_email (email) USING BTREE,
  KEY idx_phone (phone) USING BTREE,
  KEY idx_status (status) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系统管理员表'`,
	})
}
