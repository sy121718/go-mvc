// 系统附件表：统一附件管理，记录文件名、路径、大小、类型、MIME、存储方式（local/OSS等）、
// MD5、分类归属，支持全文检索文件名。
package migrations

func init() {
	register(Migration{
		Version:   "011",
		TableName: "sys_attachment",
		SQL: `CREATE TABLE sys_attachment (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  category_id bigint unsigned DEFAULT NULL COMMENT '分类ID',
  file_name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件名',
  file_path varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件路径',
  file_size bigint NOT NULL COMMENT '文件大小(字节)',
  file_type varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件类型',
  mime_type varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'MIME类型',
  storage_type varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'local' COMMENT '存储类型',
  storage_path varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '存储路径',
  url varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '访问URL',
  md5 varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件MD5',
  extra_info json DEFAULT NULL COMMENT '额外信息',
  status tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0=禁用,1=启用',
  create_by bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  update_by bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  create_time datetime(3) NOT NULL COMMENT '创建时间',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id) USING BTREE,
  KEY sys_attachment_category_id_index (category_id) USING BTREE,
  KEY idx_file_type (file_type) USING BTREE,
  KEY idx_storage_type (storage_type) USING BTREE,
  KEY idx_status (status) USING BTREE,
  KEY idx_file_size (file_size) USING BTREE,
  KEY idx_create_time (create_time) USING BTREE,
  KEY idx_update_time (update_time) USING BTREE,
  KEY idx_type_status (file_type,status) USING BTREE,
  KEY idx_status_time (status,create_time) USING BTREE,
  FULLTEXT KEY ft_file_name (file_name),
  CONSTRAINT fk_sys_attachment_category_id FOREIGN KEY (category_id) REFERENCES sys_file_category (id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系统附件表'`,
	})
}
