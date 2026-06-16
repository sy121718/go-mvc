// 邮件发送收件人表：记录每封邮件的具体收件人、发送状态（pending/sent/failed）及错误信息，关联 email_send_record。
package migrations

func init() {
	register(Migration{
		Version:   "003",
		TableName: "email_send_recipient",
		SQL: `CREATE TABLE email_send_recipient (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  record_id bigint unsigned NOT NULL COMMENT '发送记录ID',
  email varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收件人邮箱',
  name varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '收件人姓名',
  status enum('pending','sent','failed') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '发送状态',
  error_message text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  sent_time datetime(3) DEFAULT NULL COMMENT '发送时间',
  create_time datetime(3) DEFAULT NULL COMMENT '创建时间',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id) USING BTREE,
  KEY email_send_recipient_record_id_index (record_id) USING BTREE,
  KEY email_send_recipient_email_index (email) USING BTREE,
  KEY email_send_recipient_status_index (status) USING BTREE,
  CONSTRAINT fk_email_send_recipient_record_id FOREIGN KEY (record_id) REFERENCES email_send_record (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='邮件发送收件人表'`,
	})
}
