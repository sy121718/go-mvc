// 权限点定义表：记录系统全部可分配权限的信息。
// Casbin 只存授权关系（谁拥有什么权限），不存"权限列表"。
// sys_permission 作为权限点目录，提供菜单创建/编辑时选择 permission_code，
// 以及角色/用户保存权限时从 permission_code 反查 api_path/api_method。
// 与 Casbin 的关系：
//   - sys_permission 是"可分配权限清单"
//   - sys_casbin_rule 是"已授权结果"
//   - 两者通过 permission_code 关联
package migrations

func init() {
	register(Migration{
		Version:   "022",
		TableName: "sys_permission",
		SQL: `CREATE TABLE sys_permission (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  permission_code varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限编码，唯一标识，如 admin:list',
  permission_name varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限名称，如 管理员列表',
  module varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属模块，如 admin',
  api_path varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '后端接口路径，如 /api/admin/list',
  api_method varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'GET' COMMENT '请求方法 GET/POST/PUT/DELETE',
  status tinyint NOT NULL DEFAULT '1' COMMENT '状态: 0=禁用 1=启用',
  remark varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  create_by bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  create_time datetime(3) DEFAULT NULL COMMENT '创建时间',
  update_by bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY uk_permission_code (permission_code) USING BTREE,
  KEY idx_module (module) USING BTREE,
  KEY idx_status (status) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='权限点定义表：系统全部可分配权限的目录'`,
	})
}