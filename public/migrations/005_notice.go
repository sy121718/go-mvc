// 通知公告表：存储系统通知、公告、消息、警告等内容，支持置顶、弹窗、定时发布/过期、软删除。
package migrations

func init() {
	register(Migration{
		Version:   "005",
		TableName: "notice",
		SQL: `CREATE TABLE notice (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  title varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  content longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '内容',
  notice_type tinyint NOT NULL DEFAULT '1' COMMENT '类型：1系统通知、2公告、3消息、4警告',
  is_top tinyint DEFAULT '0' COMMENT '是否置顶：1是、0否',
  is_popup tinyint DEFAULT '0' COMMENT '是否弹窗提醒：1是、0否',
  publish_time datetime(3) DEFAULT NULL COMMENT '发布时间',
  expire_time datetime(3) DEFAULT NULL COMMENT '过期时间',
  read_count int DEFAULT '0' COMMENT '阅读次数',
  status tinyint DEFAULT '0' COMMENT '状态：0草稿、1已发布、2已撤回',
  create_by bigint unsigned NOT NULL COMMENT '创建人ID',
  create_time datetime(3) NOT NULL COMMENT '创建时间',
  update_by bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  update_time datetime(3) DEFAULT NULL COMMENT '更新时间',
  deleted_time datetime(3) DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (id) USING BTREE,
  KEY idx_notice_type (notice_type) USING BTREE,
  KEY idx_status (status) USING BTREE,
  KEY idx_publish_time (publish_time) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='通知公告表'`,
	})
}
