// 数据权限规则表：定义数据域（ORDERS/NOTICE/ADMIN 等）的权限规则，
// 通过 JSON 配置排除字段和条件组，实现细粒度的数据行/列级权限控制。
package migrations

func init() {
	register(Migration{
		Version:   "018",
		TableName: "sys_rule",
		SQL: `CREATE TABLE sys_rule (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  rule_name varchar(100) NOT NULL COMMENT '规则名称',
  domain varchar(50) NOT NULL COMMENT '数据域标识（ORDERS / NOTICE / ADMIN 等）',
  config json NOT NULL COMMENT '规则配置JSON（含 omit_fields + condition_groups）',
  status tinyint NOT NULL DEFAULT '1' COMMENT '状态: 0=禁用 1=启用',
  remark varchar(200) DEFAULT NULL COMMENT '备注',
  create_by bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  create_time datetime(3) DEFAULT NULL COMMENT '创建时间',
  update_by bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_domain (domain),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='数据权限规则表'`,
	})
}
