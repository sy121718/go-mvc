// 多语言文本表：国际化支持，key+lang 唯一确定一条翻译，支持错误码、UI文案、字典等多分类，可配置 HTTP 状态码。
package migrations

func init() {
	register(Migration{
		Version:   "015",
		TableName: "sys_i18n",
		SQL: `CREATE TABLE sys_i18n (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  item_key varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '键（错误码/UI文本/字典/提示等）',
  lang varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '语言（zh-CN/en-US/ja-JP等）',
  item_value text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '翻译文本',
  http_code int DEFAULT '200' COMMENT '状态码',
  category varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '分类（error/ui/dict/msg，可选）',
  remark varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注说明',
  status tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  create_time datetime(3) NOT NULL,
  update_time datetime(3) DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY uk_key_lang (item_key,lang) USING BTREE,
  KEY idx_lang (lang) USING BTREE,
  KEY idx_category (category) USING BTREE,
  KEY idx_status (status) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='多语言文本表'`,
	})
}
