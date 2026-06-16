// 规则分配表：将数据权限规则分配给角色或用户，支持按角色（target_type=1）或按用户（target_type=2）分配。
package migrations

func init() {
	register(Migration{
		Version:   "019",
		TableName: "sys_rule_assignment",
		SQL: `CREATE TABLE sys_rule_assignment (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  rule_id bigint unsigned NOT NULL COMMENT '规则ID',
  target_type tinyint NOT NULL COMMENT '目标类型: 1=角色 2=用户',
  target_id bigint unsigned NOT NULL COMMENT '目标ID（role_id 或 admin_id）',
  create_by bigint unsigned DEFAULT NULL,
  create_time datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_rule_id (rule_id),
  KEY idx_target (target_type,target_id),
  CONSTRAINT fk_rule_assign FOREIGN KEY (rule_id) REFERENCES sys_rule (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='规则分配表'`,
	})
}
