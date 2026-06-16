// 后台管理员第三方登录关联表：存储管理员与第三方平台（微信/QQ/Google 等）的绑定关系，
// 含 open_id、union_id、access_token 等，同一管理员同一平台只能绑定一次。
package migrations

func init() {
	register(Migration{
		Version:   "009",
		TableName: "sys_admin_social",
		SQL: `CREATE TABLE sys_admin_social (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  admin_id bigint unsigned DEFAULT NULL COMMENT '关联管理员ID（sys_admin.id）',
  provider_code varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '第三方平台标识',
  open_id varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '第三方平台用户唯一标识',
  union_id varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '多平台统一标识',
  access_token varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '访问令牌（加密存储）',
  expires_in int DEFAULT NULL COMMENT '令牌有效期（秒）',
  refresh_token varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '刷新令牌（加密存储）',
  bind_time datetime(3) DEFAULT NULL COMMENT '绑定时间',
  last_login_time datetime(3) DEFAULT NULL COMMENT '最后通过该平台登录的时间',
  status tinyint(1) DEFAULT '1' COMMENT '状态：1正常、0禁用',
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY uk_admin_provider (admin_id,provider_code) USING BTREE COMMENT '同一管理员同一平台只能绑定一次',
  KEY idx_provider_openid (provider_code,open_id) USING BTREE COMMENT '通过平台+openid快速查询绑定关系',
  CONSTRAINT fk_sys_admin_social_admin_id FOREIGN KEY (admin_id) REFERENCES sys_admin (id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='后台管理员第三方登录关联表'`,
	})
}
