/*
 Navicat Premium Data Transfer

 Source Server         : 本地-mysql8
 Source Server Type    : MySQL
 Source Server Version : 80012
 Source Host           : localhost:3306
 Source Schema         : base

 Target Server Type    : MySQL
 Target Server Version : 80012
 File Encoding         : 65001

 Date: 17/04/2026 14:01:46
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for email_send_recipient
-- ----------------------------
DROP TABLE IF EXISTS `email_send_recipient`;
CREATE TABLE `email_send_recipient`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `record_id` bigint(20) UNSIGNED NOT NULL COMMENT '发送记录ID',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收件人邮箱',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '收件人姓名',
  `status` enum('pending','sent','failed') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '发送状态',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '错误信息',
  `sent_time` datetime(3) NULL DEFAULT NULL COMMENT '发送时间',
  `create_time` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `email_send_recipient_record_id_index`(`record_id` ASC) USING BTREE,
  INDEX `email_send_recipient_email_index`(`email` ASC) USING BTREE,
  INDEX `email_send_recipient_status_index`(`status` ASC) USING BTREE,
  CONSTRAINT `fk_email_send_recipient_record_id` FOREIGN KEY (`record_id`) REFERENCES `email_send_record` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '邮件发送收件人表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of email_send_recipient
-- ----------------------------

-- ----------------------------
-- Table structure for email_send_record
-- ----------------------------
DROP TABLE IF EXISTS `email_send_record`;
CREATE TABLE `email_send_record`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `template_id` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '模板ID',
  `subject` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件主题',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件内容',
  `status` enum('pending','sending','success','failed') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '发送状态',
  `total_recipients` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '总收件人数',
  `success_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '成功数',
  `failed_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '失败数',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '错误信息',
  `sent_time` datetime(3) NULL DEFAULT NULL COMMENT '发送时间',
  `create_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '创建人ID',
  `create_time` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `email_send_record_template_id_index`(`template_id` ASC) USING BTREE,
  INDEX `email_send_record_status_index`(`status` ASC) USING BTREE,
  INDEX `email_send_record_sent_time_index`(`sent_time` ASC) USING BTREE,
  CONSTRAINT `fk_email_send_record_template_id` FOREIGN KEY (`template_id`) REFERENCES `email_template` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '邮件发送记录表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of email_send_record
-- ----------------------------

-- ----------------------------
-- Table structure for email_template
-- ----------------------------
DROP TABLE IF EXISTS `email_template`;
CREATE TABLE `email_template`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `template_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模板名称',
  `template_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模板编码',
  `subject` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件主题',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件内容（支持变量）',
  `template_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'html' COMMENT '模板类型: html/text',
  `variables` json NULL COMMENT '可用变量列表',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态: 0=禁用,1=启用',
  `create_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '更新人ID',
  `create_time` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `email_template_template_code_unique`(`template_code` ASC) USING BTREE,
  INDEX `email_template_status_index`(`status` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '邮件模板表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of email_template
-- ----------------------------

-- ----------------------------
-- Table structure for ip_blacklist
-- ----------------------------
DROP TABLE IF EXISTS `ip_blacklist`;
CREATE TABLE `ip_blacklist`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `ip_address` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  `error_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '错误类型',
  `ban_reason` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '拉黑原因',
  `ban_duration` int(11) NOT NULL COMMENT '封禁时长(分钟)',
  `banned_time` datetime(3) NOT NULL COMMENT '封禁时间',
  `banned_until_time` datetime(3) NOT NULL COMMENT '解封时间',
  `ban_count` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '拉黑次数',
  `operator_id` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '操作人ID',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 0=过期,1=生效',
  `create_time` datetime(3) NOT NULL COMMENT '创建时间',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ip`(`ip_address` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'IP黑名单表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of ip_blacklist
-- ----------------------------

-- ----------------------------
-- Table structure for notice
-- ----------------------------
DROP TABLE IF EXISTS `notice`;
CREATE TABLE `notice`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '内容',
  `notice_type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '类型：1系统通知、2公告、3消息、4警告',
  `is_top` tinyint(4) NULL DEFAULT 0 COMMENT '是否置顶：1是、0否',
  `is_popup` tinyint(4) NULL DEFAULT 0 COMMENT '是否弹窗提醒：1是、0否',
  `publish_time` datetime(3) NULL DEFAULT NULL COMMENT '发布时间',
  `expire_time` datetime(3) NULL DEFAULT NULL COMMENT '过期时间',
  `read_count` int(11) NULL DEFAULT 0 COMMENT '阅读次数',
  `status` tinyint(4) NULL DEFAULT 0 COMMENT '状态：0草稿、1已发布、2已撤回',
  `create_by` bigint(20) UNSIGNED NOT NULL COMMENT '创建人ID',
  `create_time` datetime(3) NOT NULL COMMENT '创建时间',
  `update_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '更新人ID',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_time` datetime(3) NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_notice_type`(`notice_type` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_publish_time`(`publish_time` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '通知公告表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of notice
-- ----------------------------

-- ----------------------------
-- Table structure for notice_read
-- ----------------------------
DROP TABLE IF EXISTS `notice_read`;
CREATE TABLE `notice_read`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `notice_id` bigint(20) UNSIGNED NOT NULL COMMENT '通知ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `read_time` datetime(3) NOT NULL COMMENT '阅读时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_notice_id`(`notice_id` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  UNIQUE INDEX `uk_notice_read_notice_user`(`notice_id` ASC, `user_id` ASC) USING BTREE,
  CONSTRAINT `fk_notice_read_notice_id` FOREIGN KEY (`notice_id`) REFERENCES `notice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '通知阅读记录表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of notice_read
-- ----------------------------

-- ----------------------------
-- Table structure for notice_target
-- ----------------------------
DROP TABLE IF EXISTS `notice_target`;
CREATE TABLE `notice_target`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `notice_id` bigint(20) UNSIGNED NOT NULL COMMENT '通知ID（关联notice.id）',
  `target_type` tinyint(4) NOT NULL COMMENT '目标类型：1全部用户、2指定角色、3指定用户',
  `target_id` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '目标ID（根据target_type关联对应表）',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_notice_id`(`notice_id` ASC) USING BTREE,
  INDEX `idx_target_type`(`target_type` ASC) USING BTREE,
  CONSTRAINT `fk_notice_target_notice_id` FOREIGN KEY (`notice_id`) REFERENCES `notice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '通知目标关联表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of notice_target
-- ----------------------------

-- ----------------------------
-- Table structure for sys_admin
-- ----------------------------
DROP TABLE IF EXISTS `sys_admin`;
CREATE TABLE `sys_admin`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID（唯一）',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录账号用户名',
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '加密密码（如bcrypt）',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '用户姓名',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '头像URL',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '手机号',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：1启用、2禁用、3密码错误封禁',
  `is_admin` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否管理员：0否、1是',
  `login_failure_count` smallint(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT '连续登录失败次数（达到阈值后临时锁定）',
  `locked_until_time` datetime(3) NULL DEFAULT NULL COMMENT '封禁至（NULL表示未封禁）',
  `last_failure_time` datetime(3) NULL DEFAULT NULL COMMENT '最后一次登录失败时间',
  `register_ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '注册IP地址',
  `register_location` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '注册地理位置（如：北京市-联通）',
  `last_login_ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '最后登录IP',
  `last_login_location` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '最后登录地理位置',
  `last_login_isp` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '最后登录网络运营商',
  `last_login_time` datetime(3) NULL DEFAULT NULL COMMENT '最后登录时间',
  `create_by` bigint(20) UNSIGNED NOT NULL COMMENT '创建人ID（0=系统创建）',
  `create_time` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint(20) UNSIGNED NOT NULL COMMENT '更新人ID（0=系统更新）',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_time` datetime(3) NULL DEFAULT NULL COMMENT '软删除时间（NULL表示未删除）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_username`(`username` ASC) USING BTREE COMMENT '用户名唯一索引',
  INDEX `idx_email`(`email` ASC) USING BTREE,
  INDEX `idx_phone`(`phone` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统管理员表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_admin
-- ----------------------------
INSERT INTO `sys_admin` VALUES (1, 'admin', '123456', '管理员', NULL, NULL, NULL, 1, 1, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0, NULL, 0, NULL, NULL);

-- ----------------------------
-- Table structure for sys_admin_sessions
-- ----------------------------
DROP TABLE IF EXISTS `sys_admin_sessions`;
CREATE TABLE `sys_admin_sessions`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '会话ID',
  `admin_id` bigint(20) UNSIGNED NOT NULL COMMENT '管理员ID（关联sys_admin.id）',
  `device_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '设备类型：desktop,mobile,tablet',
  `device_info` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '设备简要信息（如：Chrome浏览器、iPhone、iPad）',
  `jwt_token` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'JWT Token',
  `login_ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录IP',
  `login_time` datetime(3) NOT NULL COMMENT '登录时间',
  `last_active_time` datetime(3) NOT NULL COMMENT '最后活跃时间',
  `expires_time` datetime(3) NOT NULL COMMENT 'Token过期时间',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：1有效，0无效',
  `create_time` datetime(3) NOT NULL COMMENT '创建时间',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `create_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '创建人ID',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_admin_id`(`admin_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_expires_time`(`expires_time` ASC) USING BTREE,
  INDEX `idx_admin_jwt`(`admin_id` ASC, `jwt_token`(100) ASC) USING BTREE,
  CONSTRAINT `fk_sys_admin_sessions_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `sys_admin` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '管理员会话表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_admin_sessions
-- ----------------------------

-- ----------------------------
-- Table structure for sys_admin_social
-- ----------------------------
DROP TABLE IF EXISTS `sys_admin_social`;
CREATE TABLE `sys_admin_social`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `admin_id` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '关联管理员ID（sys_admin.id）',
  `provider_code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '第三方平台标识（关联social_login_providers.provider_code，如wechat/qq/google）',
  `open_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '第三方平台用户唯一标识（如微信openid、谷歌sub）',
  `union_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '多平台统一标识（如微信unionid，非必填）',
  `access_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '访问令牌（加密存储）',
  `expires_in` int(11) NULL DEFAULT NULL COMMENT '令牌有效期（秒）',
  `refresh_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '刷新令牌（加密存储）',
  `bind_time` datetime(3) NULL DEFAULT NULL COMMENT '绑定时间',
  `last_login_time` datetime(3) NULL DEFAULT NULL COMMENT '最后通过该平台登录的时间',
  `status` tinyint(1) NULL DEFAULT 1 COMMENT '状态：1正常、0禁用',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_admin_provider`(`admin_id` ASC, `provider_code` ASC) USING BTREE COMMENT '同一管理员同一平台只能绑定一次',
  INDEX `idx_provider_openid`(`provider_code` ASC, `open_id` ASC) USING BTREE COMMENT '通过平台+openid快速查询绑定关系',
  CONSTRAINT `fk_sys_admin_social_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `sys_admin` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '后台管理员第三方登录关联表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_admin_social
-- ----------------------------

-- ----------------------------
-- Table structure for sys_attachment
-- ----------------------------
DROP TABLE IF EXISTS `sys_attachment`;
CREATE TABLE `sys_attachment`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `category_id` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '分类ID',
  `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件名',
  `file_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件路径',
  `file_size` bigint(20) NOT NULL COMMENT '文件大小(字节)',
  `file_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件类型',
  `mime_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'MIME类型',
  `storage_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'local' COMMENT '存储类型',
  `storage_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '存储路径',
  `url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '访问URL',
  `md5` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '文件MD5',
  `extra_info` json NULL COMMENT '额外信息',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 0=禁用,1=启用',
  `create_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '更新人ID',
  `create_time` datetime(3) NOT NULL COMMENT '创建时间',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `sys_attachment_category_id_index`(`category_id` ASC) USING BTREE,
  INDEX `idx_file_type`(`file_type` ASC) USING BTREE,
  INDEX `idx_storage_type`(`storage_type` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `idx_file_size`(`file_size` ASC) USING BTREE,
  INDEX `idx_create_time`(`create_time` ASC) USING BTREE,
  INDEX `idx_update_time`(`update_time` ASC) USING BTREE,
  INDEX `idx_type_status`(`file_type` ASC, `status` ASC) USING BTREE,
  INDEX `idx_status_time`(`status` ASC, `create_time` ASC) USING BTREE,
  FULLTEXT INDEX `ft_file_name`(`file_name`),
  CONSTRAINT `fk_sys_attachment_category_id` FOREIGN KEY (`category_id`) REFERENCES `sys_file_category` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统附件表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_attachment
-- ----------------------------

-- ----------------------------
-- Table structure for sys_casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `sys_casbin_rule`;
CREATE TABLE `sys_casbin_rule`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `ptype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v0` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v1` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v2` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v3` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v4` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `v5` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_casbin_rule`(`ptype` ASC, `v0` ASC, `v1` ASC, `v2` ASC, `v3` ASC, `v4` ASC, `v5` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Casbin 权限策略表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_casbin_rule
-- ----------------------------

-- ----------------------------
-- Table structure for sys_config
-- ----------------------------
DROP TABLE IF EXISTS `sys_config`;
CREATE TABLE `sys_config`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `group_key` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置分组键名',
  `group_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置分组名称',
  `config_data` json NOT NULL COMMENT '配置数据(JSON格式)',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注说明',
  `status` tinyint(4) NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
  `create_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '创建人ID',
  `create_time` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '更新人ID',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_group_key`(`group_key` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统配置表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_config
-- ----------------------------
INSERT INTO `sys_config` VALUES (1, 'cache_version', '缓存版本控制', '{\"version\": 1}', '用于错误消息、i18n、配置的版本控制', 1, NULL, '2026-04-14 14:45:04.379', NULL, NULL);

-- ----------------------------
-- Table structure for sys_cron_job
-- ----------------------------
DROP TABLE IF EXISTS `sys_cron_job`;
CREATE TABLE `sys_cron_job`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `job_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务名称',
  `job_category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务分类',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '任务描述',
  `cron_expression` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Cron表达式',
  `command` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '执行命令',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 0=禁用,1=启用',
  `sort_order` int(11) NULL DEFAULT NULL COMMENT '排序',
  `last_sync_time` datetime(3) NULL DEFAULT NULL COMMENT '上次执行时间',
  `create_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '更新人ID',
  `create_time` datetime(3) NOT NULL COMMENT '创建时间',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '定时任务表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_cron_job
-- ----------------------------

-- ----------------------------
-- Table structure for sys_file_category
-- ----------------------------
DROP TABLE IF EXISTS `sys_file_category`;
CREATE TABLE `sys_file_category`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `category_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `category_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类编码',
  `parent_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级ID',
  `sort_order` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '图标',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态',
  `create_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `update_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '更新人',
  `create_time` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `sys_file_category_category_code_unique`(`category_code` ASC) USING BTREE,
  INDEX `sys_file_category_parent_id_index`(`parent_id` ASC) USING BTREE,
  INDEX `sys_file_category_status_index`(`status` ASC) USING BTREE,
  INDEX `sys_file_category_create_by_foreign`(`create_by` ASC) USING BTREE,
  INDEX `sys_file_category_update_by_foreign`(`update_by` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文件分类表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_file_category
-- ----------------------------

-- ----------------------------
-- Table structure for sys_i18n
-- ----------------------------
DROP TABLE IF EXISTS `sys_i18n`;
CREATE TABLE `sys_i18n`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `item_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '键（错误码/UI文本/字典/提示等）',
  `lang` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '语言（zh-CN/en-US/ja-JP等）',
  `item_value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '翻译文本',
  `http_code` int(11) NULL DEFAULT 200 COMMENT '状态码',
  `category` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '分类（error/ui/dict/msg，可选）',
  `remark` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注说明',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态',
  `create_time` datetime(3) NOT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_key_lang`(`item_key` ASC, `lang` ASC) USING BTREE,
  INDEX `idx_lang`(`lang` ASC) USING BTREE,
  INDEX `idx_category`(`category` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 59 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '多语言文本表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_i18n
-- ----------------------------
INSERT INTO `sys_i18n` VALUES (1, 'ErrSystemError', 'zh-CN', '系统异常', 500, 'error', '系统级错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (2, 'ErrSystemError', 'en-US', 'System error', 500, 'error', 'System error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (3, 'ErrDBQueryError', 'zh-CN', '数据库查询错误', 500, 'error', '系统级错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (4, 'ErrDBQueryError', 'en-US', 'Database query error', 500, 'error', 'System error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (5, 'ErrCacheError', 'zh-CN', '缓存错误', 500, 'error', '系统级错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (6, 'ErrCacheError', 'en-US', 'Cache error', 500, 'error', 'System error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (7, 'ErrInvalidParams', 'zh-CN', '请求参数错误', 400, 'error', '请求错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (8, 'ErrInvalidParams', 'en-US', 'Invalid request parameters', 400, 'error', 'Request error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (9, 'ErrInvalidBody', 'zh-CN', '请求体格式错误', 400, 'error', '请求错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (10, 'ErrInvalidBody', 'en-US', 'Invalid request body', 400, 'error', 'Request error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (11, 'ErrUserNotFound', 'zh-CN', '用户不存在', 404, 'error', '用户模块错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (12, 'ErrUserNotFound', 'en-US', 'User not found', 404, 'error', 'User module error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (13, 'ErrUserExists', 'zh-CN', '用户已存在', 400, 'error', '用户模块错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (14, 'ErrUserExists', 'en-US', 'User already exists', 400, 'error', 'User module error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (15, 'ErrInvalidPassword', 'zh-CN', '用户名或密码错误', 400, 'error', '用户模块错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (16, 'ErrInvalidPassword', 'en-US', 'Invalid username or password', 400, 'error', 'User module error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (17, 'ErrAdminNotFound', 'zh-CN', '管理员不存在', 404, 'error', 'Admin模块错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (18, 'ErrAdminNotFound', 'en-US', 'Admin not found', 404, 'error', 'Admin module error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (19, 'ErrAdminExists', 'zh-CN', '管理员已存在', 400, 'error', 'Admin模块错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (20, 'ErrAdminExists', 'en-US', 'Admin already exists', 400, 'error', 'Admin module error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (21, 'ErrAdminDisabled', 'zh-CN', '管理员已被禁用', 403, 'error', 'Admin模块错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (22, 'ErrAdminDisabled', 'en-US', 'Admin has been disabled', 403, 'error', 'Admin module error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (23, 'ErrUnauthorized', 'zh-CN', '未登录或登录已过期', 401, 'error', '认证错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (24, 'ErrUnauthorized', 'en-US', 'Unauthorized', 401, 'error', 'Auth error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (25, 'ErrInvalidToken', 'zh-CN', 'Token无效', 401, 'error', '认证错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (26, 'ErrInvalidToken', 'en-US', 'Invalid token', 401, 'error', 'Auth error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (27, 'ErrTokenExpired', 'zh-CN', 'Token已过期', 401, 'error', '认证错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (28, 'ErrTokenExpired', 'en-US', 'Token expired', 401, 'error', 'Auth error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (29, 'ErrPermissionDenied', 'zh-CN', '无权限访问', 403, 'error', '权限错误', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (30, 'ErrPermissionDenied', 'en-US', 'Permission denied', 403, 'error', 'Permission error', 1, '2026-04-14 16:31:04.302', NULL);
INSERT INTO `sys_i18n` VALUES (31, 'ui_admin_menu_dashboard', 'zh-CN', '仪表盘', 200, 'ui', '菜单文本', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (32, 'ui_admin_menu_dashboard', 'en-US', 'Dashboard', 200, 'ui', 'Menu text', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (33, 'ui_admin_menu_user', 'zh-CN', '用户管理', 200, 'ui', '菜单文本', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (34, 'ui_admin_menu_user', 'en-US', 'User Management', 200, 'ui', 'Menu text', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (35, 'ui_admin_btn_save', 'zh-CN', '保存', 200, 'ui', '按钮文本', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (36, 'ui_admin_btn_save', 'en-US', 'Save', 200, 'ui', 'Button text', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (37, 'ui_admin_btn_cancel', 'zh-CN', '取消', 200, 'ui', '按钮文本', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (38, 'ui_admin_btn_cancel', 'en-US', 'Cancel', 200, 'ui', 'Button text', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (39, 'ui_admin_btn_delete', 'zh-CN', '删除', 200, 'ui', '按钮文本', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (40, 'ui_admin_btn_delete', 'en-US', 'Delete', 200, 'ui', 'Button text', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (41, 'dict_admin_status_normal', 'zh-CN', '正常', 200, 'dict', '管理员状态', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (42, 'dict_admin_status_normal', 'en-US', 'Active', 200, 'dict', 'Admin status', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (43, 'dict_admin_status_disabled', 'zh-CN', '禁用', 200, 'dict', '管理员状态', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (44, 'dict_admin_status_disabled', 'en-US', 'Disabled', 200, 'dict', 'Admin status', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (45, 'dict_user_gender_male', 'zh-CN', '男', 200, 'dict', '用户性别', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (46, 'dict_user_gender_male', 'en-US', 'Male', 200, 'dict', 'User gender', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (47, 'dict_user_gender_female', 'zh-CN', '女', 200, 'dict', '用户性别', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (48, 'dict_user_gender_female', 'en-US', 'Female', 200, 'dict', 'User gender', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (49, 'msg_save_success', 'zh-CN', '保存成功', 200, 'msg', '操作提示', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (50, 'msg_save_success', 'en-US', 'Saved successfully', 200, 'msg', 'Operation message', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (51, 'msg_delete_success', 'zh-CN', '删除成功', 200, 'msg', '操作提示', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (52, 'msg_delete_success', 'en-US', 'Deleted successfully', 200, 'msg', 'Operation message', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (53, 'msg_update_success', 'zh-CN', '更新成功', 200, 'msg', '操作提示', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (54, 'msg_update_success', 'en-US', 'Updated successfully', 200, 'msg', 'Operation message', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (55, 'msg_delete_confirm', 'zh-CN', '确定要删除吗？', 200, 'msg', '确认提示', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (56, 'msg_delete_confirm', 'en-US', 'Are you sure to delete?', 200, 'msg', 'Confirm message', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (57, 'msg_operation_success', 'zh-CN', '操作成功', 200, 'msg', '操作提示', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (58, 'msg_operation_success', 'en-US', 'Operation successful', 200, 'msg', 'Operation message', 1, '2026-04-14 16:31:20.726', NULL);
INSERT INTO `sys_i18n` VALUES (59, 'ErrUploadSystemError', 'zh-CN', '上传系统异常', 500, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (60, 'ErrUploadSystemError', 'en-US', 'Upload system error', 500, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (61, 'ErrUploadNotInitialized', 'zh-CN', '上传组件未初始化', 500, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (62, 'ErrUploadNotInitialized', 'en-US', 'Upload component is not initialized', 500, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (63, 'ErrUploadProviderNotFound', 'zh-CN', '上传存储提供者不存在', 400, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (64, 'ErrUploadProviderNotFound', 'en-US', 'Upload provider not found', 400, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (65, 'ErrUploadConfigMissing', 'zh-CN', '上传配置缺失', 400, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (66, 'ErrUploadConfigMissing', 'en-US', 'Upload config is missing', 400, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (67, 'ErrUploadConfigInvalid', 'zh-CN', '上传配置无效', 400, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (68, 'ErrUploadConfigInvalid', 'en-US', 'Upload config is invalid', 400, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (69, 'ErrUploadFileEmpty', 'zh-CN', '上传文件不能为空', 400, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (70, 'ErrUploadFileEmpty', 'en-US', 'Upload file cannot be empty', 400, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (71, 'ErrUploadFileNameRequired', 'zh-CN', '上传文件名不能为空', 400, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (72, 'ErrUploadFileNameRequired', 'en-US', 'Upload file name is required', 400, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (73, 'ErrUploadWriteFailed', 'zh-CN', '上传写入失败', 500, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (74, 'ErrUploadWriteFailed', 'en-US', 'Upload write failed', 500, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (75, 'ErrUploadRequestFailed', 'zh-CN', '上传请求失败', 502, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (76, 'ErrUploadRequestFailed', 'en-US', 'Upload request failed', 502, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (77, 'ErrUploadTokenFailed', 'zh-CN', '上传签名生成失败', 500, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (78, 'ErrUploadTokenFailed', 'en-US', 'Upload token/sign generation failed', 500, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (79, 'ErrUploadResponseInvalid', 'zh-CN', '上传响应无效', 502, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');
INSERT INTO `sys_i18n` VALUES (80, 'ErrUploadResponseInvalid', 'en-US', 'Upload response is invalid', 502, 'error', 'upload error code', 1, '2026-04-17 11:35:30.734', '2026-04-17 11:35:30.734');

-- ----------------------------
-- Table structure for sys_logs
-- ----------------------------
DROP TABLE IF EXISTS `sys_logs`;
CREATE TABLE `sys_logs`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `log_type` tinyint(4) NOT NULL COMMENT '日志类型：1=登录 2=操作',
  `admin_id` bigint(20) UNSIGNED NOT NULL COMMENT '管理员ID',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  `status` tinyint(4) NOT NULL COMMENT '状态：1=成功 0=失败',
  `api_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求接口路径',
  `http_method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求方法：GET/POST',
  `operation` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作名称：登录/登出/创建用户等',
  `device_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '设备类型：desktop/mobile/tablet',
  `location` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '地理位置',
  `user_agent` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '浏览器信息',
  `detail` json NULL COMMENT '日志详情',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`, `created_at`) USING BTREE,
  INDEX `idx_admin_time`(`admin_id` ASC, `created_at` ASC) USING BTREE,
  INDEX `idx_type_admin_time`(`log_type` ASC, `admin_id` ASC, `created_at` ASC) USING BTREE,
  INDEX `idx_ip_time`(`ip` ASC, `created_at` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统日志表（分区）' ROW_FORMAT = DYNAMIC PARTITION BY RANGE COLUMNS (`created_at`)
PARTITIONS 4
(PARTITION `p202603` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p202604` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p202605` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 ,
PARTITION `p_future` ENGINE = InnoDB MAX_ROWS = 0 MIN_ROWS = 0 )
;

-- ----------------------------
-- Records of sys_logs
-- ----------------------------

-- ----------------------------
-- Table structure for sys_menus
-- ----------------------------
DROP TABLE IF EXISTS `sys_menus`;
CREATE TABLE `sys_menus`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `menu_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单编码，唯一标识',
  `permission_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '权限编码，对应 Casbin obj 字段',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '菜单标题',
  `parent_id` bigint(20) NULL DEFAULT 0 COMMENT '父级ID，0为顶级',
  `type` tinyint(4) NOT NULL DEFAULT 2 COMMENT '类型：1=目录 2=菜单 3=按钮 4=iframe 5=外链',
  `path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '前端路由路径',
  `component` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '前端组件路径',
  `external_url` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '外链/iframe地址',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '菜单图标',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：0=禁用 1=启用',
  `is_hidden` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否隐藏：0=显示 1=隐藏',
  `is_public` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否公开：0=需权限 1=无需权限',
  `is_system` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否系统内置：0=否 1=是（不可删除）',
  `sort_order` int(11) NULL DEFAULT 0 COMMENT '排序',
  `remark` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注',
  `create_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `create_time` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '更新人',
  `update_time` datetime(3) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_time` datetime(3) NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_menu_code`(`menu_code` ASC) USING BTREE,
  INDEX `idx_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统菜单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_menus
-- ----------------------------

-- ----------------------------
-- Procedure structure for sys_logs_add_partition
-- ----------------------------
DROP PROCEDURE IF EXISTS `sys_logs_add_partition`;
delimiter ;;
CREATE PROCEDURE `sys_logs_add_partition`()
BEGIN
    DECLARE next_month_start VARCHAR(20);
    DECLARE next_next_month_start VARCHAR(20);
    
    SET next_month_start = DATE_FORMAT(DATE_ADD(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-01');
    SET next_next_month_start = DATE_FORMAT(DATE_ADD(next_month_start, INTERVAL 1 MONTH), '%Y-%m-01');
    
    SET @sql = CONCAT(
        'ALTER TABLE `sys_logs` REORGANIZE PARTITION `p_future` INTO (',
        'PARTITION `p', DATE_FORMAT(next_month_start, '%%Y%%m'), '` VALUES LESS THAN (''', next_next_month_start, '''),',
        'PARTITION `p_future` VALUES LESS THAN MAXVALUE',
        ')'
    );
    
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END
;;
delimiter ;

-- ----------------------------
-- Event structure for evt_sys_logs_add_partition
-- ----------------------------
DROP EVENT IF EXISTS `evt_sys_logs_add_partition`;
delimiter ;;
CREATE EVENT `evt_sys_logs_add_partition`
ON SCHEDULE
EVERY '1' MONTH STARTS '2026-04-01 00:00:00'
ON COMPLETION PRESERVE
COMMENT '每月自动为 sys_logs 添加新分区'
DO CALL sys_logs_add_partition()
;;
delimiter ;

SET FOREIGN_KEY_CHECKS = 1;
