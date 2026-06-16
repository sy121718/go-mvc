// Casbin 权限策略表：Casbin RBAC 引擎的策略存储表，动态管理用户/角色与资源/操作的权限关系。
package migrations

func init() {
	register(Migration{
		Version:   "012",
		TableName: "sys_casbin_rule",
		SQL: `CREATE TABLE sys_casbin_rule (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  ptype varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  v0 varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  v1 varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  v2 varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  v3 varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  v4 varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  v5 varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY idx_sys_casbin_rule (ptype,v0,v1,v2,v3,v4,v5) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='Casbin 权限策略表'`,
	})
}
