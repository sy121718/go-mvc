-- MySQL dump 10.13  Distrib 8.4.9, for Linux (x86_64)
--
-- Host: localhost    Database: base
-- ------------------------------------------------------
-- Server version	8.4.9

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `email_send_recipient`
--

DROP TABLE IF EXISTS `email_send_recipient`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `email_send_recipient` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `record_id` bigint unsigned NOT NULL COMMENT '发送记录ID',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收件人邮箱',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '收件人姓名',
  `status` enum('pending','sent','failed') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '发送状态',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  `sent_time` datetime(3) DEFAULT NULL COMMENT '发送时间',
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `email_send_recipient_record_id_index` (`record_id`) USING BTREE,
  KEY `email_send_recipient_email_index` (`email`) USING BTREE,
  KEY `email_send_recipient_status_index` (`status`) USING BTREE,
  CONSTRAINT `fk_email_send_recipient_record_id` FOREIGN KEY (`record_id`) REFERENCES `email_send_record` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='邮件发送收件人表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `email_send_recipient`
--

LOCK TABLES `email_send_recipient` WRITE;
/*!40000 ALTER TABLE `email_send_recipient` DISABLE KEYS */;
/*!40000 ALTER TABLE `email_send_recipient` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `email_send_record`
--

DROP TABLE IF EXISTS `email_send_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `email_send_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `template_id` bigint unsigned DEFAULT NULL COMMENT '模板ID',
  `subject` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件主题',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件内容',
  `status` enum('pending','sending','success','failed') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '发送状态',
  `total_recipients` int unsigned NOT NULL DEFAULT '0' COMMENT '总收件人数',
  `success_count` int unsigned NOT NULL DEFAULT '0' COMMENT '成功数',
  `failed_count` int unsigned NOT NULL DEFAULT '0' COMMENT '失败数',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  `sent_time` datetime(3) DEFAULT NULL COMMENT '发送时间',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `email_send_record_template_id_index` (`template_id`) USING BTREE,
  KEY `email_send_record_status_index` (`status`) USING BTREE,
  KEY `email_send_record_sent_time_index` (`sent_time`) USING BTREE,
  CONSTRAINT `fk_email_send_record_template_id` FOREIGN KEY (`template_id`) REFERENCES `email_template` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='邮件发送记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `email_send_record`
--

LOCK TABLES `email_send_record` WRITE;
/*!40000 ALTER TABLE `email_send_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `email_send_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `email_template`
--

DROP TABLE IF EXISTS `email_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `email_template` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `template_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模板名称',
  `template_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模板编码',
  `subject` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件主题',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件内容（支持变量）',
  `template_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'html' COMMENT '模板类型: html/text',
  `variables` json DEFAULT NULL COMMENT '可用变量列表',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态: 0=禁用,1=启用',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `email_template_template_code_unique` (`template_code`) USING BTREE,
  KEY `email_template_status_index` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='邮件模板表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `email_template`
--

LOCK TABLES `email_template` WRITE;
/*!40000 ALTER TABLE `email_template` DISABLE KEYS */;
/*!40000 ALTER TABLE `email_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ip_blacklist`
--

DROP TABLE IF EXISTS `ip_blacklist`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ip_blacklist` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `ip_address` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  `error_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '错误类型',
  `ban_reason` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '拉黑原因',
  `ban_duration` int NOT NULL COMMENT '封禁时长(分钟)',
  `banned_time` datetime(3) NOT NULL COMMENT '封禁时间',
  `banned_until_time` datetime(3) NOT NULL COMMENT '解封时间',
  `ban_count` int unsigned NOT NULL DEFAULT '1' COMMENT '拉黑次数',
  `operator_id` bigint unsigned DEFAULT NULL COMMENT '操作人ID',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0=过期,1=生效',
  `create_time` datetime(3) NOT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_ip` (`ip_address`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='IP黑名单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ip_blacklist`
--

LOCK TABLES `ip_blacklist` WRITE;
/*!40000 ALTER TABLE `ip_blacklist` DISABLE KEYS */;
/*!40000 ALTER TABLE `ip_blacklist` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notice`
--

DROP TABLE IF EXISTS `notice`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notice` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '内容',
  `notice_type` tinyint NOT NULL DEFAULT '1' COMMENT '类型：1系统通知、2公告、3消息、4警告',
  `is_top` tinyint DEFAULT '0' COMMENT '是否置顶：1是、0否',
  `is_popup` tinyint DEFAULT '0' COMMENT '是否弹窗提醒：1是、0否',
  `publish_time` datetime(3) DEFAULT NULL COMMENT '发布时间',
  `expire_time` datetime(3) DEFAULT NULL COMMENT '过期时间',
  `read_count` int DEFAULT '0' COMMENT '阅读次数',
  `status` tinyint DEFAULT '0' COMMENT '状态：0草稿、1已发布、2已撤回',
  `create_by` bigint unsigned NOT NULL COMMENT '创建人ID',
  `create_time` datetime(3) NOT NULL COMMENT '创建时间',
  `update_by` bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_notice_type` (`notice_type`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE,
  KEY `idx_publish_time` (`publish_time`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='通知公告表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notice`
--

LOCK TABLES `notice` WRITE;
/*!40000 ALTER TABLE `notice` DISABLE KEYS */;
/*!40000 ALTER TABLE `notice` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notice_read`
--

DROP TABLE IF EXISTS `notice_read`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notice_read` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `notice_id` bigint unsigned NOT NULL COMMENT '通知ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `read_time` datetime(3) NOT NULL COMMENT '阅读时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_notice_read_notice_user` (`notice_id`,`user_id`) USING BTREE,
  KEY `idx_notice_id` (`notice_id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  CONSTRAINT `fk_notice_read_notice_id` FOREIGN KEY (`notice_id`) REFERENCES `notice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='通知阅读记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notice_read`
--

LOCK TABLES `notice_read` WRITE;
/*!40000 ALTER TABLE `notice_read` DISABLE KEYS */;
/*!40000 ALTER TABLE `notice_read` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notice_target`
--

DROP TABLE IF EXISTS `notice_target`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notice_target` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `notice_id` bigint unsigned NOT NULL COMMENT '通知ID（关联notice.id）',
  `target_type` tinyint NOT NULL COMMENT '目标类型：1全部用户、2指定角色、3指定用户',
  `target_id` bigint unsigned DEFAULT NULL COMMENT '目标ID（根据target_type关联对应表）',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_notice_id` (`notice_id`) USING BTREE,
  KEY `idx_target_type` (`target_type`) USING BTREE,
  CONSTRAINT `fk_notice_target_notice_id` FOREIGN KEY (`notice_id`) REFERENCES `notice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='通知目标关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notice_target`
--

LOCK TABLES `notice_target` WRITE;
/*!40000 ALTER TABLE `notice_target` DISABLE KEYS */;
/*!40000 ALTER TABLE `notice_target` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_admin`
--

DROP TABLE IF EXISTS `sys_admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_admin` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID（唯一）',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录账号用户名',
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '加密密码（如bcrypt）',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户姓名',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '头像URL',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '手机号',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1启用、2禁用、3密码错误封禁',
  `is_admin` tinyint NOT NULL DEFAULT '0' COMMENT '是否管理员：0否、1是',
  `login_failure_count` smallint unsigned NOT NULL DEFAULT '0' COMMENT '连续登录失败次数（达到阈值后临时锁定）',
  `locked_until_time` datetime(3) DEFAULT NULL COMMENT '封禁至（NULL表示未封禁）',
  `metadata` json DEFAULT NULL COMMENT '扩展元数据',
  `last_failure_time` datetime(3) DEFAULT NULL COMMENT '最后一次登录失败时间',
  `register_ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '注册IP地址',
  `register_location` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '注册地理位置（如：北京市-联通）',
  `last_login_ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后登录IP',
  `last_login_location` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后登录地理位置',
  `last_login_isp` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后登录网络运营商',
  `last_login_time` datetime(3) DEFAULT NULL COMMENT '最后登录时间',
  `create_by` bigint unsigned NOT NULL COMMENT '创建人ID（0=系统创建）',
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint unsigned NOT NULL COMMENT '更新人ID（0=系统更新）',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_username` (`username`) USING BTREE COMMENT '用户名唯一索引',
  KEY `idx_email` (`email`) USING BTREE,
  KEY `idx_phone` (`phone`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系统管理员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_admin`
--

LOCK TABLES `sys_admin` WRITE;
/*!40000 ALTER TABLE `sys_admin` DISABLE KEYS */;
INSERT INTO `sys_admin` VALUES (1,'admin','$2a$10$Xk2cyVAGlTtfBtzTOiD5ae4FYHzLSz9lR1ggGx3r.I7CLCer.Hg4y','admin','https://avatars.githubusercontent.com/u/52823142','1217189608@qq.com',NULL,1,1,0,NULL,NULL,NULL,NULL,NULL,'::1',NULL,NULL,'2026-06-02 14:37:54.758',0,NULL,0,NULL,NULL);
/*!40000 ALTER TABLE `sys_admin` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_admin_social`
--

DROP TABLE IF EXISTS `sys_admin_social`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_admin_social` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `admin_id` bigint unsigned DEFAULT NULL COMMENT '关联管理员ID（sys_admin.id）',
  `provider_code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '第三方平台标识（关联social_login_providers.provider_code，如wechat/qq/google）',
  `open_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '第三方平台用户唯一标识（如微信openid、谷歌sub）',
  `union_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '多平台统一标识（如微信unionid，非必填）',
  `access_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '访问令牌（加密存储）',
  `expires_in` int DEFAULT NULL COMMENT '令牌有效期（秒）',
  `refresh_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '刷新令牌（加密存储）',
  `bind_time` datetime(3) DEFAULT NULL COMMENT '绑定时间',
  `last_login_time` datetime(3) DEFAULT NULL COMMENT '最后通过该平台登录的时间',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态：1正常、0禁用',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_admin_provider` (`admin_id`,`provider_code`) USING BTREE COMMENT '同一管理员同一平台只能绑定一次',
  KEY `idx_provider_openid` (`provider_code`,`open_id`) USING BTREE COMMENT '通过平台+openid快速查询绑定关系',
  CONSTRAINT `fk_sys_admin_social_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `sys_admin` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='后台管理员第三方登录关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_admin_social`
--

LOCK TABLES `sys_admin_social` WRITE;
/*!40000 ALTER TABLE `sys_admin_social` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_admin_social` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_attachment`
--

DROP TABLE IF EXISTS `sys_attachment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_attachment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `category_id` bigint unsigned DEFAULT NULL COMMENT '分类ID',
  `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件名',
  `file_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件路径',
  `file_size` bigint NOT NULL COMMENT '文件大小(字节)',
  `file_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件类型',
  `mime_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'MIME类型',
  `storage_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'local' COMMENT '存储类型',
  `storage_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '存储路径',
  `url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '访问URL',
  `md5` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件MD5',
  `extra_info` json DEFAULT NULL COMMENT '额外信息',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0=禁用,1=启用',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  `create_time` datetime(3) NOT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `sys_attachment_category_id_index` (`category_id`) USING BTREE,
  KEY `idx_file_type` (`file_type`) USING BTREE,
  KEY `idx_storage_type` (`storage_type`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE,
  KEY `idx_file_size` (`file_size`) USING BTREE,
  KEY `idx_create_time` (`create_time`) USING BTREE,
  KEY `idx_update_time` (`update_time`) USING BTREE,
  KEY `idx_type_status` (`file_type`,`status`) USING BTREE,
  KEY `idx_status_time` (`status`,`create_time`) USING BTREE,
  FULLTEXT KEY `ft_file_name` (`file_name`),
  CONSTRAINT `fk_sys_attachment_category_id` FOREIGN KEY (`category_id`) REFERENCES `sys_file_category` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系统附件表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_attachment`
--

LOCK TABLES `sys_attachment` WRITE;
/*!40000 ALTER TABLE `sys_attachment` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_attachment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_casbin_rule`
--

DROP TABLE IF EXISTS `sys_casbin_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_casbin_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `ptype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v0` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v1` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v2` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v3` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v4` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v5` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_sys_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='Casbin 权限策略表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_casbin_rule`
--

LOCK TABLES `sys_casbin_rule` WRITE;
/*!40000 ALTER TABLE `sys_casbin_rule` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_casbin_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_config`
--

DROP TABLE IF EXISTS `sys_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_config` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `group_key` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置分组键名',
  `group_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置分组名称',
  `config_data` json NOT NULL COMMENT '配置数据(JSON格式)',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注说明',
  `status` tinyint DEFAULT '1' COMMENT '状态：1-启用，0-禁用',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_group_key` (`group_key`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系统配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_config`
--

LOCK TABLES `sys_config` WRITE;
/*!40000 ALTER TABLE `sys_config` DISABLE KEYS */;
INSERT INTO `sys_config` VALUES (1,'cache_version','缓存版本控制','{\"version\": 1}','用于错误消息、i18n、配置的版本控制',1,NULL,'2026-04-14 14:45:04.379',NULL,NULL);
/*!40000 ALTER TABLE `sys_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_cron_job`
--

DROP TABLE IF EXISTS `sys_cron_job`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_cron_job` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `job_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务名称',
  `job_category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务分类',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '任务描述',
  `cron_expression` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Cron表达式',
  `command` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '执行命令',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0=禁用,1=启用',
  `sort_order` int DEFAULT NULL COMMENT '排序',
  `last_sync_time` datetime(3) DEFAULT NULL COMMENT '上次执行时间',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  `create_time` datetime(3) NOT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='定时任务表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_cron_job`
--

LOCK TABLES `sys_cron_job` WRITE;
/*!40000 ALTER TABLE `sys_cron_job` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_cron_job` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_file_category`
--

DROP TABLE IF EXISTS `sys_file_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_file_category` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `category_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `category_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类编码',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父级ID',
  `sort_order` int unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '图标',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建人',
  `update_by` bigint unsigned DEFAULT NULL COMMENT '更新人',
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `sys_file_category_category_code_unique` (`category_code`) USING BTREE,
  KEY `sys_file_category_parent_id_index` (`parent_id`) USING BTREE,
  KEY `sys_file_category_status_index` (`status`) USING BTREE,
  KEY `sys_file_category_create_by_foreign` (`create_by`) USING BTREE,
  KEY `sys_file_category_update_by_foreign` (`update_by`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='文件分类表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_file_category`
--

LOCK TABLES `sys_file_category` WRITE;
/*!40000 ALTER TABLE `sys_file_category` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_file_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_i18n`
--

DROP TABLE IF EXISTS `sys_i18n`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_i18n` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `item_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '键（错误码/UI文本/字典/提示等）',
  `lang` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '语言（zh-CN/en-US/ja-JP等）',
  `item_value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '翻译文本',
  `http_code` int DEFAULT '200' COMMENT '状态码',
  `category` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '分类（error/ui/dict/msg，可选）',
  `remark` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注说明',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `create_time` datetime(3) NOT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_key_lang` (`item_key`,`lang`) USING BTREE,
  KEY `idx_lang` (`lang`) USING BTREE,
  KEY `idx_category` (`category`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=81 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='多语言文本表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_i18n`
--

LOCK TABLES `sys_i18n` WRITE;
/*!40000 ALTER TABLE `sys_i18n` DISABLE KEYS */;
INSERT INTO `sys_i18n` VALUES (1,'ErrSystemError','zh-CN','系统异常',500,'error','系统级错误',1,'2026-04-14 16:31:04.302',NULL),(2,'ErrSystemError','en-US','System error',500,'error','System error',1,'2026-04-14 16:31:04.302',NULL),(3,'ErrDBQueryError','zh-CN','数据库查询错误',500,'error','系统级错误',1,'2026-04-14 16:31:04.302',NULL),(4,'ErrDBQueryError','en-US','Database query error',500,'error','System error',1,'2026-04-14 16:31:04.302',NULL),(5,'ErrCacheError','zh-CN','缓存错误',500,'error','系统级错误',1,'2026-04-14 16:31:04.302',NULL),(6,'ErrCacheError','en-US','Cache error',500,'error','System error',1,'2026-04-14 16:31:04.302',NULL),(7,'ErrInvalidParams','zh-CN','请求参数错误',400,'error','请求错误',1,'2026-04-14 16:31:04.302',NULL),(8,'ErrInvalidParams','en-US','Invalid request parameters',400,'error','Request error',1,'2026-04-14 16:31:04.302',NULL),(9,'ErrInvalidBody','zh-CN','请求体格式错误',400,'error','请求错误',1,'2026-04-14 16:31:04.302',NULL),(10,'ErrInvalidBody','en-US','Invalid request body',400,'error','Request error',1,'2026-04-14 16:31:04.302',NULL),(11,'ErrUserNotFound','zh-CN','用户不存在',404,'error','用户模块错误',1,'2026-04-14 16:31:04.302',NULL),(12,'ErrUserNotFound','en-US','User not found',404,'error','User module error',1,'2026-04-14 16:31:04.302',NULL),(13,'ErrUserExists','zh-CN','用户已存在',400,'error','用户模块错误',1,'2026-04-14 16:31:04.302',NULL),(14,'ErrUserExists','en-US','User already exists',400,'error','User module error',1,'2026-04-14 16:31:04.302',NULL),(15,'ErrInvalidPassword','zh-CN','用户名或密码错误',400,'error','用户模块错误',1,'2026-04-14 16:31:04.302',NULL),(16,'ErrInvalidPassword','en-US','Invalid username or password',400,'error','User module error',1,'2026-04-14 16:31:04.302',NULL),(17,'ErrAdminNotFound','zh-CN','管理员不存在',404,'error','Admin模块错误',1,'2026-04-14 16:31:04.302',NULL),(18,'ErrAdminNotFound','en-US','Admin not found',404,'error','Admin module error',1,'2026-04-14 16:31:04.302',NULL),(19,'ErrAdminExists','zh-CN','管理员已存在',400,'error','Admin模块错误',1,'2026-04-14 16:31:04.302',NULL),(20,'ErrAdminExists','en-US','Admin already exists',400,'error','Admin module error',1,'2026-04-14 16:31:04.302',NULL),(21,'ErrAdminDisabled','zh-CN','管理员已被禁用',403,'error','Admin模块错误',1,'2026-04-14 16:31:04.302',NULL),(22,'ErrAdminDisabled','en-US','Admin has been disabled',403,'error','Admin module error',1,'2026-04-14 16:31:04.302',NULL),(23,'ErrUnauthorized','zh-CN','未登录或登录已过期',401,'error','认证错误',1,'2026-04-14 16:31:04.302',NULL),(24,'ErrUnauthorized','en-US','Unauthorized',401,'error','Auth error',1,'2026-04-14 16:31:04.302',NULL),(25,'ErrInvalidToken','zh-CN','Token无效',401,'error','认证错误',1,'2026-04-14 16:31:04.302',NULL),(26,'ErrInvalidToken','en-US','Invalid token',401,'error','Auth error',1,'2026-04-14 16:31:04.302',NULL),(27,'ErrTokenExpired','zh-CN','Token已过期',401,'error','认证错误',1,'2026-04-14 16:31:04.302',NULL),(28,'ErrTokenExpired','en-US','Token expired',401,'error','Auth error',1,'2026-04-14 16:31:04.302',NULL),(29,'ErrPermissionDenied','zh-CN','无权限访问',403,'error','权限错误',1,'2026-04-14 16:31:04.302',NULL),(30,'ErrPermissionDenied','en-US','Permission denied',403,'error','Permission error',1,'2026-04-14 16:31:04.302',NULL),(31,'ui_admin_menu_dashboard','zh-CN','仪表盘',200,'ui','菜单文本',1,'2026-04-14 16:31:20.726',NULL),(32,'ui_admin_menu_dashboard','en-US','Dashboard',200,'ui','Menu text',1,'2026-04-14 16:31:20.726',NULL),(33,'ui_admin_menu_user','zh-CN','用户管理',200,'ui','菜单文本',1,'2026-04-14 16:31:20.726',NULL),(34,'ui_admin_menu_user','en-US','User Management',200,'ui','Menu text',1,'2026-04-14 16:31:20.726',NULL),(35,'ui_admin_btn_save','zh-CN','保存',200,'ui','按钮文本',1,'2026-04-14 16:31:20.726',NULL),(36,'ui_admin_btn_save','en-US','Save',200,'ui','Button text',1,'2026-04-14 16:31:20.726',NULL),(37,'ui_admin_btn_cancel','zh-CN','取消',200,'ui','按钮文本',1,'2026-04-14 16:31:20.726',NULL),(38,'ui_admin_btn_cancel','en-US','Cancel',200,'ui','Button text',1,'2026-04-14 16:31:20.726',NULL),(39,'ui_admin_btn_delete','zh-CN','删除',200,'ui','按钮文本',1,'2026-04-14 16:31:20.726',NULL),(40,'ui_admin_btn_delete','en-US','Delete',200,'ui','Button text',1,'2026-04-14 16:31:20.726',NULL),(41,'dict_admin_status_normal','zh-CN','正常',200,'dict','管理员状态',1,'2026-04-14 16:31:20.726',NULL),(42,'dict_admin_status_normal','en-US','Active',200,'dict','Admin status',1,'2026-04-14 16:31:20.726',NULL),(43,'dict_admin_status_disabled','zh-CN','禁用',200,'dict','管理员状态',1,'2026-04-14 16:31:20.726',NULL),(44,'dict_admin_status_disabled','en-US','Disabled',200,'dict','Admin status',1,'2026-04-14 16:31:20.726',NULL),(45,'dict_user_gender_male','zh-CN','男',200,'dict','用户性别',1,'2026-04-14 16:31:20.726',NULL),(46,'dict_user_gender_male','en-US','Male',200,'dict','User gender',1,'2026-04-14 16:31:20.726',NULL),(47,'dict_user_gender_female','zh-CN','女',200,'dict','用户性别',1,'2026-04-14 16:31:20.726',NULL),(48,'dict_user_gender_female','en-US','Female',200,'dict','User gender',1,'2026-04-14 16:31:20.726',NULL),(49,'msg_save_success','zh-CN','保存成功',200,'msg','操作提示',1,'2026-04-14 16:31:20.726',NULL),(50,'msg_save_success','en-US','Saved successfully',200,'msg','Operation message',1,'2026-04-14 16:31:20.726',NULL),(51,'msg_delete_success','zh-CN','删除成功',200,'msg','操作提示',1,'2026-04-14 16:31:20.726',NULL),(52,'msg_delete_success','en-US','Deleted successfully',200,'msg','Operation message',1,'2026-04-14 16:31:20.726',NULL),(53,'msg_update_success','zh-CN','更新成功',200,'msg','操作提示',1,'2026-04-14 16:31:20.726',NULL),(54,'msg_update_success','en-US','Updated successfully',200,'msg','Operation message',1,'2026-04-14 16:31:20.726',NULL),(55,'msg_delete_confirm','zh-CN','确定要删除吗？',200,'msg','确认提示',1,'2026-04-14 16:31:20.726',NULL),(56,'msg_delete_confirm','en-US','Are you sure to delete?',200,'msg','Confirm message',1,'2026-04-14 16:31:20.726',NULL),(57,'msg_operation_success','zh-CN','操作成功',200,'msg','操作提示',1,'2026-04-14 16:31:20.726',NULL),(58,'msg_operation_success','en-US','Operation successful',200,'msg','Operation message',1,'2026-04-14 16:31:20.726',NULL),(59,'ErrUploadSystemError','zh-CN','上传系统异常',500,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(60,'ErrUploadSystemError','en-US','Upload system error',500,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(61,'ErrUploadNotInitialized','zh-CN','上传组件未初始化',500,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(62,'ErrUploadNotInitialized','en-US','Upload component is not initialized',500,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(63,'ErrUploadProviderNotFound','zh-CN','上传存储提供者不存在',400,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(64,'ErrUploadProviderNotFound','en-US','Upload provider not found',400,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(65,'ErrUploadConfigMissing','zh-CN','上传配置缺失',400,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(66,'ErrUploadConfigMissing','en-US','Upload config is missing',400,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(67,'ErrUploadConfigInvalid','zh-CN','上传配置无效',400,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(68,'ErrUploadConfigInvalid','en-US','Upload config is invalid',400,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(69,'ErrUploadFileEmpty','zh-CN','上传文件不能为空',400,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(70,'ErrUploadFileEmpty','en-US','Upload file cannot be empty',400,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(71,'ErrUploadFileNameRequired','zh-CN','上传文件名不能为空',400,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(72,'ErrUploadFileNameRequired','en-US','Upload file name is required',400,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(73,'ErrUploadWriteFailed','zh-CN','上传写入失败',500,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(74,'ErrUploadWriteFailed','en-US','Upload write failed',500,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(75,'ErrUploadRequestFailed','zh-CN','上传请求失败',502,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(76,'ErrUploadRequestFailed','en-US','Upload request failed',502,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(77,'ErrUploadTokenFailed','zh-CN','上传签名生成失败',500,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(78,'ErrUploadTokenFailed','en-US','Upload token/sign generation failed',500,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(79,'ErrUploadResponseInvalid','zh-CN','上传响应无效',502,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734'),(80,'ErrUploadResponseInvalid','en-US','Upload response is invalid',502,'error','upload error code',1,'2026-04-17 11:35:30.734','2026-04-17 11:35:30.734');
/*!40000 ALTER TABLE `sys_i18n` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_logs`
--

DROP TABLE IF EXISTS `sys_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `log_type` tinyint NOT NULL COMMENT '日志类型：1=登录 2=操作',
  `admin_id` bigint unsigned NOT NULL COMMENT '管理员ID',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  `status` tinyint NOT NULL COMMENT '状态：1=成功 0=失败',
  `api_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求接口路径',
  `http_method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求方法：GET/POST',
  `operation` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作名称：登录/登出/创建用户等',
  `device_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '设备类型：desktop/mobile/tablet',
  `location` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '地理位置',
  `user_agent` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '浏览器信息',
  `detail` json DEFAULT NULL COMMENT '日志详情',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`created_at`) USING BTREE,
  KEY `idx_admin_time` (`admin_id`,`created_at`) USING BTREE,
  KEY `idx_type_admin_time` (`log_type`,`admin_id`,`created_at`) USING BTREE,
  KEY `idx_ip_time` (`ip`,`created_at`) USING BTREE,
  KEY `idx_created_at` (`created_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系统日志表（分区）'
/*!50500 PARTITION BY RANGE  COLUMNS(created_at)
(PARTITION p202603 VALUES LESS THAN ('2026-04-01 00:00:00') ENGINE = InnoDB,
 PARTITION p202604 VALUES LESS THAN ('2026-05-01 00:00:00') ENGINE = InnoDB,
 PARTITION p202605 VALUES LESS THAN ('2026-06-01 00:00:00') ENGINE = InnoDB,
 PARTITION p202606 VALUES LESS THAN ('2026-07-01 00:00:00') ENGINE = InnoDB,
 PARTITION `p%Y%m` VALUES LESS THAN ('2026-08-01') ENGINE = InnoDB,
 PARTITION p_future VALUES LESS THAN (MAXVALUE) ENGINE = InnoDB) */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_logs`
--

LOCK TABLES `sys_logs` WRITE;
/*!40000 ALTER TABLE `sys_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_menus`
--

DROP TABLE IF EXISTS `sys_menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_menus` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `menu_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单编码，唯一标识',
  `permission_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '权限编码，对应 Casbin obj 字段',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '菜单标题',
  `parent_id` bigint DEFAULT '0' COMMENT '父级ID，0为顶级',
  `type` tinyint NOT NULL DEFAULT '2' COMMENT '类型：1=目录 2=菜单 3=按钮 4=iframe 5=外链',
  `path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '前端路由路径',
  `component` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '前端组件路径',
  `external_url` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '外链/iframe地址',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '菜单图标',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0=禁用 1=启用',
  `is_hidden` tinyint NOT NULL DEFAULT '0' COMMENT '是否隐藏：0=显示 1=隐藏',
  `is_public` tinyint NOT NULL DEFAULT '0' COMMENT '是否公开：0=需权限 1=无需权限',
  `is_system` tinyint NOT NULL DEFAULT '0' COMMENT '是否系统内置：0=否 1=是（不可删除）',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `remark` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建人',
  `create_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_by` bigint unsigned DEFAULT NULL COMMENT '更新人',
  `update_time` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_time` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_menu_code` (`menu_code`) USING BTREE,
  KEY `idx_parent_id` (`parent_id`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系统菜单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_menus`
--

LOCK TABLES `sys_menus` WRITE;
/*!40000 ALTER TABLE `sys_menus` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_rule`
--

DROP TABLE IF EXISTS `sys_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `rule_name` varchar(100) NOT NULL COMMENT '规则名称',
  `domain` varchar(50) NOT NULL COMMENT '数据域标识（ORDERS / NOTICE / ADMIN 等）',
  `config` json NOT NULL COMMENT '规则配置JSON（含 omit_fields + condition_groups）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态: 0=禁用 1=启用',
  `remark` varchar(200) DEFAULT NULL COMMENT '备注',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_domain` (`domain`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='数据权限规则表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_rule`
--

LOCK TABLES `sys_rule` WRITE;
/*!40000 ALTER TABLE `sys_rule` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_rule_assignment`
--

DROP TABLE IF EXISTS `sys_rule_assignment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_rule_assignment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `rule_id` bigint unsigned NOT NULL COMMENT '规则ID',
  `target_type` tinyint NOT NULL COMMENT '目标类型: 1=角色 2=用户',
  `target_id` bigint unsigned NOT NULL COMMENT '目标ID（role_id 或 admin_id）',
  `create_by` bigint unsigned DEFAULT NULL,
  `create_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_rule_id` (`rule_id`),
  KEY `idx_target` (`target_type`,`target_id`),
  CONSTRAINT `fk_rule_assign` FOREIGN KEY (`rule_id`) REFERENCES `sys_rule` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='规则分配表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_rule_assignment`
--

LOCK TABLES `sys_rule_assignment` WRITE;
/*!40000 ALTER TABLE `sys_rule_assignment` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_rule_assignment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'base'
--
/*!50106 SET @save_time_zone= @@TIME_ZONE */ ;
/*!50106 DROP EVENT IF EXISTS `evt_sys_logs_add_partition` */;
DELIMITER ;;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;;
/*!50003 SET character_set_client  = utf8mb4 */ ;;
/*!50003 SET character_set_results = utf8mb4 */ ;;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;;
/*!50003 SET @saved_time_zone      = @@time_zone */ ;;
/*!50003 SET time_zone             = 'SYSTEM' */ ;;
/*!50106 CREATE*/ /*!50117 DEFINER=`root`@`localhost`*/ /*!50106 EVENT `evt_sys_logs_add_partition` ON SCHEDULE EVERY 1 MONTH STARTS '2026-04-01 00:00:00' ON COMPLETION PRESERVE ENABLE COMMENT '每月自动为 sys_logs 添加新分区' DO CALL sys_logs_add_partition() */ ;;
/*!50003 SET time_zone             = @saved_time_zone */ ;;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;;
/*!50003 SET character_set_client  = @saved_cs_client */ ;;
/*!50003 SET character_set_results = @saved_cs_results */ ;;
/*!50003 SET collation_connection  = @saved_col_connection */ ;;
DELIMITER ;
/*!50106 SET TIME_ZONE= @save_time_zone */ ;

--
-- Dumping routines for database 'base'
--
/*!50003 DROP PROCEDURE IF EXISTS `sys_logs_add_partition` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `sys_logs_add_partition`()
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
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-06-02 15:29:18
