// 通知阅读记录表：记录用户对通知的已读状态，同一用户对同一通知仅记录一次（唯一约束）。
package migrations

func init() {
	register(Migration{
		Version:   "006",
		TableName: "notice_read",
		SQL: `CREATE TABLE notice_read (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  notice_id bigint unsigned NOT NULL COMMENT '通知ID',
  user_id bigint unsigned NOT NULL COMMENT '用户ID',
  read_time datetime(3) NOT NULL COMMENT '阅读时间',
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY uk_notice_read_notice_user (notice_id,user_id) USING BTREE,
  KEY idx_notice_id (notice_id) USING BTREE,
  KEY idx_user_id (user_id) USING BTREE,
  CONSTRAINT fk_notice_read_notice_id FOREIGN KEY (notice_id) REFERENCES notice (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='通知阅读记录表'`,
	})
}
