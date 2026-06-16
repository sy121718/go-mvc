// 邮件发送记录表：记录每次邮件发送的 subject、content、状态（pending/sending/success/failed）、
// 收件人统计（总数/成功/失败），关联 email_template。
package migrations

func init() {
	register(Migration{
		Version:   "002",
		TableName: "email_send_record",
		SQL: `CREATE TABLE email_send_record (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  template_id bigint unsigned DEFAULT NULL COMMENT '模板ID',
  subject varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件主题',
  content text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件内容',
  status enum('pending','sending','success','failed') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '发送状态',
  total_recipients int unsigned NOT NULL DEFAULT '0' COMMENT '总收件人数',
  success_count int unsigned NOT NULL DEFAULT '0' COMMENT '成功数',
  failed_count int unsigned NOT NULL DEFAULT '0' COMMENT '失败数',
  error_message text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  sent_time datetime(3) DEFAULT NULL COMMENT '发送时间',
  create_by bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  create_time datetime(3) DEFAULT NULL COMMENT '创建时间',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id) USING BTREE,
  KEY email_send_record_template_id_index (template_id) USING BTREE,
  KEY email_send_record_status_index (status) USING BTREE,
  KEY email_send_record_sent_time_index (sent_time) USING BTREE,
  CONSTRAINT fk_email_send_record_template_id FOREIGN KEY (template_id) REFERENCES email_template (id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='邮件发送记录表'`,
	})
}
