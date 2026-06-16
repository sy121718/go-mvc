// 系统日志表（分区表）：记录管理员登录和操作日志，按 created_at 列范围分区，
// 每月自动生成新分区（配合存储过程 sys_logs_add_partition），支持按管理员、IP、时间等维度查询。
package migrations

func init() {
	register(Migration{
		Version:   "016",
		TableName: "sys_logs",
		SQL: `CREATE TABLE sys_logs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  log_type tinyint NOT NULL COMMENT '日志类型：1=登录 2=操作',
  admin_id bigint unsigned NOT NULL COMMENT '管理员ID',
  ip varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  status tinyint NOT NULL COMMENT '状态：1=成功 0=失败',
  api_path varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求接口路径',
  http_method varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求方法：GET/POST',
  operation varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作名称：登录/登出/创建用户等',
  device_type varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '设备类型：desktop/mobile/tablet',
  location varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '地理位置',
  user_agent varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '浏览器信息',
  detail json DEFAULT NULL COMMENT '日志详情',
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id,created_at) USING BTREE,
  KEY idx_admin_time (admin_id,created_at) USING BTREE,
  KEY idx_type_admin_time (log_type,admin_id,created_at) USING BTREE,
  KEY idx_ip_time (ip,created_at) USING BTREE,
  KEY idx_created_at (created_at) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系统日志表（分区）'
PARTITION BY RANGE COLUMNS(created_at) (
  PARTITION p202603 VALUES LESS THAN ('2026-04-01 00:00:00') ENGINE = InnoDB,
  PARTITION p202604 VALUES LESS THAN ('2026-05-01 00:00:00') ENGINE = InnoDB,
  PARTITION p202605 VALUES LESS THAN ('2026-06-01 00:00:00') ENGINE = InnoDB,
  PARTITION p202606 VALUES LESS THAN ('2026-07-01 00:00:00') ENGINE = InnoDB,
  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = InnoDB
)`,
	})
}
