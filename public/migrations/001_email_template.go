// 邮件模板表：存储邮件模板名称、编码、主题、内容（支持变量）、类型（html/text）等。
package migrations

func init() {
	register(Migration{
		Version:   "001",
		TableName: "email_template",
		SQL: `CREATE TABLE email_template (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  template_name varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模板名称',
  template_code varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模板编码',
  subject varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件主题',
  content text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件内容（支持变量）',
  template_type varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'html' COMMENT '模板类型: html/text',
  variables json DEFAULT NULL COMMENT '可用变量列表',
  status tinyint NOT NULL DEFAULT '1' COMMENT '状态: 0=禁用,1=启用',
  create_by bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  update_by bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  create_time datetime(3) DEFAULT NULL COMMENT '创建时间',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY email_template_template_code_unique (template_code) USING BTREE,
  KEY email_template_status_index (status) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='邮件模板表'`,
	})
}
