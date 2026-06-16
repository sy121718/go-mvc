// 系统配置表：分组存储系统配置数据（JSON 格式），如缓存版本控制等，通过 group_key 唯一标识。
package migrations

func init() {
	register(Migration{
		Version:   "013",
		TableName: "sys_config",
		SQL: `CREATE TABLE sys_config (
  id int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  group_key varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置分组键名',
  group_name varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置分组名称',
  config_data json NOT NULL COMMENT '配置数据(JSON格式)',
  remark varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注说明',
  status tinyint DEFAULT '1' COMMENT '状态：1-启用，0-禁用',
  create_by bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  create_time datetime(3) DEFAULT NULL COMMENT '创建时间',
  update_by bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY uk_group_key (group_key) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系统配置表'`,
	})
}
