// 文件分类表：文件/附件的分类管理，支持多级父子分类结构、排序、图标等。
package migrations

func init() {
	register(Migration{
		Version:   "010",
		TableName: "sys_file_category",
		SQL: `CREATE TABLE sys_file_category (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  category_name varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  category_code varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类编码',
  parent_id bigint unsigned NOT NULL DEFAULT '0' COMMENT '父级ID',
  sort_order int unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  icon varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '图标',
  status tinyint NOT NULL DEFAULT '1' COMMENT '状态',
  create_by bigint unsigned DEFAULT NULL COMMENT '创建人',
  update_by bigint unsigned DEFAULT NULL COMMENT '更新人',
  create_time datetime(3) DEFAULT NULL COMMENT '创建时间',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY sys_file_category_category_code_unique (category_code) USING BTREE,
  KEY sys_file_category_parent_id_index (parent_id) USING BTREE,
  KEY sys_file_category_status_index (status) USING BTREE,
  KEY sys_file_category_create_by_foreign (create_by) USING BTREE,
  KEY sys_file_category_update_by_foreign (update_by) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='文件分类表'`,
	})
}
