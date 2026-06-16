// 系统菜单表：管理后台菜单树，支持目录/菜单/按钮/iframe/外链五种类型。
// 菜单即"权限的可视化"：type=2,3 的 permission_code 是权限编码（关联 Casbin，唯一），
// type=1,4 不参与权限（permission_code 为 NULL），前端用 path/id 标识。
// 含路由、组件路径、图标、排序、隐藏/公开/系统内置等属性，支持软删除。
package migrations

func init() {
	register(Migration{
		Version:   "017",
		TableName: "sys_menus",
		SQL: `CREATE TABLE sys_menus (
  id bigint NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  permission_code varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '权限编码，对应 Casbin 关联键；type=2,3 必填且唯一',
  title varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '菜单标题',
  parent_id bigint DEFAULT '0' COMMENT '父级ID，0为顶级',
  type tinyint NOT NULL DEFAULT '2' COMMENT '类型：1=目录 2=菜单 3=按钮 4=iframe 5=外链',
  path varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '前端路由路径',
  component varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '前端组件路径',
  external_url varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '外链/iframe地址',
  icon varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '菜单图标',
  status tinyint NOT NULL DEFAULT '1' COMMENT '状态：0=禁用 1=启用',
  is_hidden tinyint NOT NULL DEFAULT '0' COMMENT '是否隐藏：0=显示 1=隐藏',
  is_public tinyint NOT NULL DEFAULT '0' COMMENT '是否公开：0=需权限 1=无需权限',
  is_system tinyint NOT NULL DEFAULT '0' COMMENT '是否系统内置：0=否 1=是（不可删除）',
  sort_order int DEFAULT '0' COMMENT '排序',
  remark varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  create_by bigint unsigned DEFAULT NULL COMMENT '创建人',
  create_time datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  update_by bigint unsigned DEFAULT NULL COMMENT '更新人',
  update_time datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  deleted_time datetime(3) DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY uk_permission_code (permission_code) USING BTREE,
  KEY idx_parent_id (parent_id) USING BTREE,
  KEY idx_status (status) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系统菜单表'`,
	})
}
