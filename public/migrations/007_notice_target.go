// 通知目标关联表：定义通知的推送范围，支持全部用户/指定角色/指定用户三种目标类型。
package migrations

func init() {
	register(Migration{
		Version:   "007",
		TableName: "notice_target",
		SQL: `CREATE TABLE notice_target (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  notice_id bigint unsigned NOT NULL COMMENT '通知ID（关联notice.id）',
  target_type tinyint NOT NULL COMMENT '目标类型：1全部用户、2指定角色、3指定用户',
  target_id bigint unsigned DEFAULT NULL COMMENT '目标ID（根据target_type关联对应表）',
  PRIMARY KEY (id) USING BTREE,
  KEY idx_notice_id (notice_id) USING BTREE,
  KEY idx_target_type (target_type) USING BTREE,
  CONSTRAINT fk_notice_target_notice_id FOREIGN KEY (notice_id) REFERENCES notice (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='通知目标关联表'`,
	})
}
