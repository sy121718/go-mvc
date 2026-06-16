// IP黑名单表：记录被封禁的 IP 地址、错误类型、封禁时长、解封时间、拉黑次数等，用于登录/访问安全控制。
package migrations

func init() {
	register(Migration{
		Version:   "004",
		TableName: "ip_blacklist",
		SQL: `CREATE TABLE ip_blacklist (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  ip_address varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  error_type varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '错误类型',
  ban_reason varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '拉黑原因',
  ban_duration int NOT NULL COMMENT '封禁时长(分钟)',
  banned_time datetime(3) NOT NULL COMMENT '封禁时间',
  banned_until_time datetime(3) NOT NULL COMMENT '解封时间',
  ban_count int unsigned NOT NULL DEFAULT '1' COMMENT '拉黑次数',
  operator_id bigint unsigned DEFAULT NULL COMMENT '操作人ID',
  remark varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  status tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0=过期,1=生效',
  create_time datetime(3) NOT NULL COMMENT '创建时间',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id) USING BTREE,
  KEY idx_ip (ip_address) USING BTREE,
  KEY idx_status (status) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='IP黑名单表'`,
	})
}
