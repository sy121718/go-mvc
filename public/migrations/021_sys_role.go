// 角色表：仅存角色元信息（名称、状态、排序、备注等）。
// 角色作为"权限组"使用：
//   - 角色含哪些权限（permission_code）由 Casbin p 策略承载：p, role_code, path, method, code
//   - 用户↔角色关系由 Casbin g 映射承载：g, user_id, role_code
// 因此不需要 sys_admin_role 关联表。sys_role 通过 role_code 与 Casbin 的 sub 关联。
package migrations

func init() {
	register(Migration{
		Version:   "021",
		TableName: "sys_role",
		SQL: `CREATE TABLE sys_role (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  role_code varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色编码，与 Casbin sub 对应的关联键',
  role_name varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称（前端展示）',
  status tinyint NOT NULL DEFAULT '1' COMMENT '状态: 0=禁用 1=启用',
  sort_order int NOT NULL DEFAULT '0' COMMENT '排序',
  remark varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  create_by bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  create_time datetime(3) DEFAULT NULL COMMENT '创建时间',
  update_by bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY uk_role_code (role_code) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='角色表（仅存元信息，权限关系在 Casbin）'`,
	})
}
