// 定时任务表：管理系统定时任务，含 Cron 表达式、执行命令、分类、排序、上次执行时间等。
package migrations

func init() {
	register(Migration{
		Version:   "014",
		TableName: "sys_cron_job",
		SQL: `CREATE TABLE sys_cron_job (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  job_name varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务名称',
  job_category varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务分类',
  description text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '任务描述',
  cron_expression varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Cron表达式',
  command varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '执行命令',
  status tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0=禁用,1=启用',
  sort_order int DEFAULT NULL COMMENT '排序',
  last_sync_time datetime(3) DEFAULT NULL COMMENT '上次执行时间',
  create_by bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  update_by bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  create_time datetime(3) NOT NULL COMMENT '创建时间',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id) USING BTREE,
  KEY idx_status (status) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='定时任务表'`,
	})
}
