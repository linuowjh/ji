-- 云念纪念馆数据库初始化脚本

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `yun_nian_memorial` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `yun_nian_memorial`;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` varchar(36) NOT NULL,
  `wechat_openid` varchar(100) NOT NULL,
  `wechat_unionid` varchar(100) DEFAULT NULL,
  `nickname` varchar(50) DEFAULT NULL,
  `avatar_url` varchar(255) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `status` tinyint DEFAULT '1' COMMENT '1:正常 0:禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_wechat_openid` (`wechat_openid`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建纪念馆表
CREATE TABLE IF NOT EXISTS `memorials` (
  `id` varchar(36) NOT NULL,
  `creator_id` varchar(36) NOT NULL,
  `deceased_name` varchar(50) NOT NULL,
  `birth_date` date DEFAULT NULL,
  `death_date` date DEFAULT NULL,
  `biography` text,
  `avatar_url` varchar(255) DEFAULT NULL,
  `theme_style` varchar(50) DEFAULT 'traditional',
  `tombstone_style` varchar(50) DEFAULT 'marble',
  `epitaph` text,
  `privacy_level` tinyint DEFAULT '1' COMMENT '1:家族可见 2:私密',
  `status` tinyint DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_memorials_creator_id` (`creator_id`),
  KEY `idx_memorials_deleted_at` (`deleted_at`),
  KEY `idx_memorials_creator_status` (`creator_id`,`status`),
  KEY `idx_memorials_privacy_status` (`privacy_level`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建祭扫记录表
CREATE TABLE IF NOT EXISTS `worship_records` (
  `id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `worship_type` varchar(20) NOT NULL COMMENT 'flower|candle|incense|tribute|prayer',
  `content` json DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_worship_records_memorial_id` (`memorial_id`),
  KEY `idx_worship_records_user_id` (`user_id`),
  KEY `idx_worship_records_deleted_at` (`deleted_at`),
  KEY `idx_worship_memorial_time` (`memorial_id`,`created_at`),
  KEY `idx_worship_user_time` (`user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建家族表
CREATE TABLE IF NOT EXISTS `families` (
  `id` varchar(36) NOT NULL,
  `name` varchar(100) NOT NULL,
  `creator_id` varchar(36) NOT NULL,
  `description` text,
  `invite_code` varchar(20) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_families_invite_code` (`invite_code`),
  KEY `idx_families_creator_id` (`creator_id`),
  KEY `idx_families_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建家族成员表
CREATE TABLE IF NOT EXISTS `family_members` (
  `id` varchar(36) NOT NULL,
  `family_id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `role` varchar(20) DEFAULT 'member' COMMENT 'admin|member',
  `joined_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_family_user` (`family_id`,`user_id`),
  KEY `idx_family_members_family_id` (`family_id`),
  KEY `idx_family_members_user_id` (`user_id`),
  KEY `idx_family_members_family` (`family_id`,`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建媒体文件表
CREATE TABLE IF NOT EXISTS `media_files` (
  `id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) NOT NULL,
  `file_type` varchar(20) NOT NULL COMMENT 'image|video|audio',
  `file_url` varchar(255) NOT NULL,
  `file_name` varchar(255) DEFAULT NULL,
  `file_size` bigint DEFAULT NULL,
  `description` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_media_files_memorial_id` (`memorial_id`),
  KEY `idx_media_files_deleted_at` (`deleted_at`),
  KEY `idx_media_files_type` (`file_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建祈福表
CREATE TABLE IF NOT EXISTS `prayers` (
  `id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `content` text NOT NULL,
  `is_public` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_prayers_memorial_id` (`memorial_id`),
  KEY `idx_prayers_user_id` (`user_id`),
  KEY `idx_prayers_deleted_at` (`deleted_at`),
  KEY `idx_prayers_memorial_public` (`memorial_id`,`is_public`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建留言表
CREATE TABLE IF NOT EXISTS `messages` (
  `id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `message_type` varchar(20) NOT NULL COMMENT 'text|audio|video',
  `content` text,
  `media_url` varchar(255) DEFAULT NULL,
  `duration` int DEFAULT NULL COMMENT '音频/视频时长(秒)',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_messages_memorial_id` (`memorial_id`),
  KEY `idx_messages_user_id` (`user_id`),
  KEY `idx_messages_deleted_at` (`deleted_at`),
  KEY `idx_messages_memorial_time` (`memorial_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建纪念日提醒表
CREATE TABLE IF NOT EXISTS `memorial_reminders` (
  `id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) NOT NULL,
  `reminder_type` varchar(20) NOT NULL COMMENT 'birthday|death_anniversary|festival',
  `reminder_date` date NOT NULL,
  `title` varchar(100) DEFAULT NULL,
  `content` text,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_memorial_reminders_memorial_id` (`memorial_id`),
  KEY `idx_memorial_reminders_deleted_at` (`deleted_at`),
  KEY `idx_memorial_reminders_date` (`reminder_date`,`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建纪念馆家族关联表
CREATE TABLE IF NOT EXISTS `memorial_families` (
  `id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) NOT NULL,
  `family_id` varchar(36) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_memorial_family` (`memorial_id`,`family_id`),
  KEY `idx_memorial_families_memorial_id` (`memorial_id`),
  KEY `idx_memorial_families_family_id` (`family_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建访客记录表
CREATE TABLE IF NOT EXISTS `visitor_records` (
  `id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) NOT NULL,
  `visitor_id` varchar(36) NOT NULL,
  `visit_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `ip_address` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_visitor_records_memorial_id` (`memorial_id`),
  KEY `idx_visitor_records_visitor_id` (`visitor_id`),
  KEY `idx_visitor_records_time` (`visit_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 添加外键约束
ALTER TABLE `memorials` ADD CONSTRAINT `fk_memorials_creator` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
ALTER TABLE `worship_records` ADD CONSTRAINT `fk_worship_records_memorial` FOREIGN KEY (`memorial_id`) REFERENCES `memorials` (`id`) ON DELETE CASCADE;
ALTER TABLE `worship_records` ADD CONSTRAINT `fk_worship_records_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
ALTER TABLE `families` ADD CONSTRAINT `fk_families_creator` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
ALTER TABLE `family_members` ADD CONSTRAINT `fk_family_members_family` FOREIGN KEY (`family_id`) REFERENCES `families` (`id`) ON DELETE CASCADE;
ALTER TABLE `family_members` ADD CONSTRAINT `fk_family_members_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
ALTER TABLE `media_files` ADD CONSTRAINT `fk_media_files_memorial` FOREIGN KEY (`memorial_id`) REFERENCES `memorials` (`id`) ON DELETE CASCADE;
ALTER TABLE `prayers` ADD CONSTRAINT `fk_prayers_memorial` FOREIGN KEY (`memorial_id`) REFERENCES `memorials` (`id`) ON DELETE CASCADE;
ALTER TABLE `prayers` ADD CONSTRAINT `fk_prayers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
ALTER TABLE `messages` ADD CONSTRAINT `fk_messages_memorial` FOREIGN KEY (`memorial_id`) REFERENCES `memorials` (`id`) ON DELETE CASCADE;
ALTER TABLE `messages` ADD CONSTRAINT `fk_messages_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
ALTER TABLE `memorial_reminders` ADD CONSTRAINT `fk_memorial_reminders_memorial` FOREIGN KEY (`memorial_id`) REFERENCES `memorials` (`id`) ON DELETE CASCADE;
ALTER TABLE `memorial_families` ADD CONSTRAINT `fk_memorial_families_memorial` FOREIGN KEY (`memorial_id`) REFERENCES `memorials` (`id`) ON DELETE CASCADE;
ALTER TABLE `memorial_families` ADD CONSTRAINT `fk_memorial_families_family` FOREIGN KEY (`family_id`) REFERENCES `families` (`id`) ON DELETE CASCADE;
ALTER TABLE `visitor_records` ADD CONSTRAINT `fk_visitor_records_memorial` FOREIGN KEY (`memorial_id`) REFERENCES `memorials` (`id`) ON DELETE CASCADE;
ALTER TABLE `visitor_records` ADD CONSTRAINT `fk_visitor_records_visitor` FOREIGN KEY (`visitor_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

-- 插入测试数据
INSERT INTO `users` (`id`, `wechat_openid`, `nickname`, `avatar_url`, `status`) VALUES
('test-user-1', 'test_openid_1', '张三', 'https://example.com/avatar1.jpg', 1),
('test-user-2', 'test_openid_2', '李四', 'https://example.com/avatar2.jpg', 1),
('test-user-3', 'test_openid_3', '王五', 'https://example.com/avatar3.jpg', 1);

INSERT INTO `memorials` (`id`, `creator_id`, `deceased_name`, `birth_date`, `death_date`, `biography`, `theme_style`, `privacy_level`) VALUES
('test-memorial-1', 'test-user-1', '张老爷子', '1950-01-01', '2020-12-31', '张老爷子是一位慈祥的长者，一生勤劳善良，深受家人和邻里的爱戴。他热爱生活，关爱家人，是我们永远的榜样。', 'traditional', 1),
('test-memorial-2', 'test-user-2', '李奶奶', '1955-03-15', '2021-06-20', '李奶奶是一位温柔的母亲，她用自己的爱温暖着整个家庭。她的笑容永远留在我们心中。', 'elegant', 1);

INSERT INTO `families` (`id`, `name`, `creator_id`, `description`, `invite_code`) VALUES
('test-family-1', '张氏家族', 'test-user-1', '张氏家族纪念圈，传承家族情感，共同缅怀先人。', 'ZHANG001'),
('test-family-2', '李氏家族', 'test-user-2', '李氏家族纪念圈，让爱跨越时空。', 'LI002');

INSERT INTO `family_members` (`id`, `family_id`, `user_id`, `role`) VALUES
('test-member-1', 'test-family-1', 'test-user-1', 'admin'),
('test-member-2', 'test-family-1', 'test-user-3', 'member'),
('test-member-3', 'test-family-2', 'test-user-2', 'admin');

INSERT INTO `memorial_families` (`id`, `memorial_id`, `family_id`) VALUES
('test-mf-1', 'test-memorial-1', 'test-family-1'),
('test-mf-2', 'test-memorial-2', 'test-family-2');

INSERT INTO `worship_records` (`id`, `memorial_id`, `user_id`, `worship_type`, `content`) VALUES
('test-worship-1', 'test-memorial-1', 'test-user-1', 'flower', '{"flower_type": "菊花", "count": 3, "message": "爷爷，我们想您了"}'),
('test-worship-2', 'test-memorial-1', 'test-user-3', 'candle', '{"candle_type": "红烛", "duration": 24, "message": "为爷爷点亮心灯"}'),
('test-worship-3', 'test-memorial-2', 'test-user-2', 'incense', '{"incense_count": 3, "message": "妈妈，愿您在天堂安好"}');

INSERT INTO `prayers` (`id`, `memorial_id`, `user_id`, `content`, `is_public`) VALUES
('test-prayer-1', 'test-memorial-1', 'test-user-1', '爷爷，愿您在天堂安好，我们会好好生活，不让您担心。', 1),
('test-prayer-2', 'test-memorial-2', 'test-user-2', '妈妈，谢谢您给我们的爱，我们永远爱您。', 1);

INSERT INTO `messages` (`id`, `memorial_id`, `user_id`, `message_type`, `content`) VALUES
('test-message-1', 'test-memorial-1', 'test-user-1', 'text', '爷爷，今天是您的生日，我们全家都在想念您。您教给我们的做人道理，我们会一直记在心里。'),
('test-message-2', 'test-memorial-2', 'test-user-2', 'text', '妈妈，今天路过您最喜欢的花园，看到满园的花开，就想起了您的笑容。');

INSERT INTO `memorial_reminders` (`id`, `memorial_id`, `reminder_type`, `reminder_date`, `title`, `content`, `is_active`) VALUES
('test-reminder-1', 'test-memorial-1', 'birthday', '2024-01-01', '张老爷子生日', '今天是张老爷子的生日，让我们一起为他献花祈福。', 1),
('test-reminder-2', 'test-memorial-1', 'death_anniversary', '2024-12-31', '张老爷子忌日', '今天是张老爷子的忌日，愿他在天堂安好。', 1),
('test-reminder-3', 'test-memorial-2', 'birthday', '2024-03-15', '李奶奶生日', '今天是李奶奶的生日，让我们一起缅怀她的慈爱。', 1);

SET FOREIGN_KEY_CHECKS = 1;

--
 创建系统配置表
CREATE TABLE IF NOT EXISTS `system_configs` (
  `id` varchar(36) NOT NULL,
  `config_key` varchar(100) NOT NULL,
  `config_value` text,
  `config_type` varchar(50) NOT NULL COMMENT 'festival|template|system',
  `description` text,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_system_configs_key` (`config_key`),
  KEY `idx_system_configs_type` (`config_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建祭扫节日配置表
CREATE TABLE IF NOT EXISTS `festival_configs` (
  `id` varchar(36) NOT NULL,
  `name` varchar(100) NOT NULL,
  `festival_date` varchar(5) NOT NULL COMMENT 'MM-DD格式',
  `description` text,
  `reminder_days` int DEFAULT '3' COMMENT '提前几天提醒',
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_festival_configs_date` (`festival_date`),
  KEY `idx_festival_configs_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建模板配置表
CREATE TABLE IF NOT EXISTS `template_configs` (
  `id` varchar(36) NOT NULL,
  `template_type` varchar(50) NOT NULL COMMENT 'theme|tombstone|prayer',
  `template_name` varchar(100) NOT NULL,
  `template_data` json,
  `preview_url` varchar(255) DEFAULT NULL,
  `is_premium` tinyint(1) DEFAULT '0',
  `sort_order` int DEFAULT '0',
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_template_configs_type` (`template_type`),
  KEY `idx_template_configs_sort` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建数据备份记录表
CREATE TABLE IF NOT EXISTS `data_backups` (
  `id` varchar(36) NOT NULL,
  `backup_type` varchar(50) NOT NULL COMMENT 'full|incremental|user',
  `backup_path` varchar(500) NOT NULL,
  `file_size` bigint DEFAULT NULL,
  `status` varchar(20) DEFAULT 'pending' COMMENT 'pending|processing|completed|failed',
  `error_message` text,
  `created_by` varchar(36) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `completed_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_data_backups_status` (`status`),
  KEY `idx_data_backups_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建系统日志表
CREATE TABLE IF NOT EXISTS `system_logs` (
  `id` varchar(36) NOT NULL,
  `log_level` varchar(20) NOT NULL COMMENT 'info|warning|error|critical',
  `log_type` varchar(50) NOT NULL COMMENT 'admin|system|security|api',
  `user_id` varchar(36) DEFAULT NULL,
  `action` varchar(200) DEFAULT NULL,
  `details` json,
  `ip_address` varchar(45) DEFAULT NULL,
  `user_agent` varchar(500) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_system_logs_level` (`log_level`),
  KEY `idx_system_logs_type` (`log_type`),
  KEY `idx_system_logs_user_id` (`user_id`),
  KEY `idx_system_logs_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建系统监控指标表
CREATE TABLE IF NOT EXISTS `system_monitors` (
  `id` varchar(36) NOT NULL,
  `metric_type` varchar(50) NOT NULL COMMENT 'cpu|memory|disk|api|database',
  `metric_value` double DEFAULT NULL,
  `metric_unit` varchar(20) DEFAULT NULL,
  `additional_info` json,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_system_monitors_type` (`metric_type`),
  KEY `idx_system_monitors_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认祭扫节日配置
INSERT INTO `festival_configs` (`id`, `name`, `festival_date`, `description`, `reminder_days`, `is_active`) VALUES
('festival-qingming', '清明节', '04-05', '清明节是中国传统的祭祖节日', 3, 1),
('festival-zhongyuan', '中元节', '08-15', '中元节（农历七月十五）是祭祀祖先的重要节日', 3, 1),
('festival-hanyi', '寒衣节', '10-01', '寒衣节（农历十月初一）是祭祀祖先、送寒衣的节日', 3, 1),
('festival-chuxi', '除夕', '12-31', '除夕是农历年的最后一天，祭祖迎新', 3, 1);

-- 插入默认模板配置
INSERT INTO `template_configs` (`id`, `template_type`, `template_name`, `template_data`, `is_premium`, `sort_order`, `is_active`) VALUES
('theme-traditional', 'theme', '中式传统', '{"background": "traditional-bg.jpg", "color_scheme": "warm", "font": "serif"}', 0, 1, 1),
('theme-elegant', 'theme', '简约素雅', '{"background": "elegant-bg.jpg", "color_scheme": "neutral", "font": "sans-serif"}', 0, 2, 1),
('theme-nature', 'theme', '自然清新', '{"background": "nature-bg.jpg", "color_scheme": "green", "font": "sans-serif"}', 0, 3, 1),
('tombstone-marble', 'tombstone', '大理石', '{"material": "marble", "color": "white", "style": "classic"}', 0, 1, 1),
('tombstone-granite', 'tombstone', '花岗岩', '{"material": "granite", "color": "gray", "style": "modern"}', 0, 2, 1);

-- 插入默认系统配置
INSERT INTO `system_configs` (`id`, `config_key`, `config_value`, `config_type`, `description`, `is_active`) VALUES
('sys-max-memorial', 'max_memorial_per_user', '10', 'system', '每个用户最多可创建的纪念馆数量', 1),
('sys-max-upload', 'max_upload_size', '10485760', 'system', '最大上传文件大小（字节）', 1),
('sys-auto-backup', 'enable_auto_backup', 'true', 'system', '是否启用自动备份', 1),
('sys-backup-interval', 'backup_interval_hours', '24', 'system', '自动备份间隔（小时）', 1);

SET FOREIGN_KEY_CHECKS = 1;


-- 创建高级套餐表
CREATE TABLE IF NOT EXISTS `premium_packages` (
  `id` varchar(36) NOT NULL,
  `package_name` varchar(100) NOT NULL,
  `package_type` varchar(50) NOT NULL COMMENT 'memorial|service|storage',
  `description` text,
  `features` json COMMENT '功能列表',
  `price` decimal(10,2) NOT NULL,
  `duration` int DEFAULT '365' COMMENT '有效期（天）',
  `storage_size` bigint DEFAULT NULL COMMENT '存储空间（字节）',
  `is_active` tinyint(1) DEFAULT '1',
  `sort_order` int DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_premium_packages_type` (`package_type`),
  KEY `idx_premium_packages_active` (`is_active`),
  KEY `idx_premium_packages_sort` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建用户订阅表
CREATE TABLE IF NOT EXISTS `user_subscriptions` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `package_id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) DEFAULT NULL,
  `status` varchar(20) DEFAULT 'active' COMMENT 'active|expired|cancelled',
  `start_date` timestamp NOT NULL,
  `end_date` timestamp NOT NULL,
  `auto_renew` tinyint(1) DEFAULT '0',
  `payment_amount` decimal(10,2) DEFAULT NULL,
  `payment_method` varchar(20) DEFAULT NULL COMMENT 'wechat|alipay',
  `transaction_id` varchar(100) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `cancelled_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_subscriptions_user_id` (`user_id`),
  KEY `idx_user_subscriptions_memorial_id` (`memorial_id`),
  KEY `idx_user_subscriptions_status` (`status`),
  KEY `idx_user_subscriptions_end_date` (`end_date`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`package_id`) REFERENCES `premium_packages` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建纪念馆升级记录表
CREATE TABLE IF NOT EXISTS `memorial_upgrades` (
  `id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) NOT NULL,
  `subscription_id` varchar(36) NOT NULL,
  `upgrade_type` varchar(50) NOT NULL COMMENT 'theme|tombstone|storage|feature',
  `upgrade_data` json,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_memorial_upgrades_memorial_id` (`memorial_id`),
  KEY `idx_memorial_upgrades_subscription_id` (`subscription_id`),
  FOREIGN KEY (`memorial_id`) REFERENCES `memorials` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`subscription_id`) REFERENCES `user_subscriptions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建定制模板表
CREATE TABLE IF NOT EXISTS `custom_templates` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) DEFAULT NULL,
  `template_type` varchar(50) NOT NULL COMMENT 'theme|tombstone|layout',
  `template_name` varchar(100) NOT NULL,
  `template_data` json,
  `preview_url` varchar(255) DEFAULT NULL,
  `status` varchar(20) DEFAULT 'draft' COMMENT 'draft|active|archived',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_custom_templates_user_id` (`user_id`),
  KEY `idx_custom_templates_memorial_id` (`memorial_id`),
  KEY `idx_custom_templates_type` (`template_type`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建存储使用情况表
CREATE TABLE IF NOT EXISTS `storage_usages` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `used_space` bigint DEFAULT '0' COMMENT '已使用空间（字节）',
  `total_space` bigint DEFAULT '104857600' COMMENT '总空间（字节）',
  `file_count` int DEFAULT '0',
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_storage_usages_user_id` (`user_id`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建支付订单表
CREATE TABLE IF NOT EXISTS `payment_orders` (
  `id` varchar(36) NOT NULL,
  `order_no` varchar(50) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `package_id` varchar(36) NOT NULL,
  `order_type` varchar(20) NOT NULL COMMENT 'subscription|upgrade|renewal',
  `amount` decimal(10,2) NOT NULL,
  `payment_method` varchar(20) DEFAULT NULL COMMENT 'wechat|alipay',
  `payment_status` varchar(20) DEFAULT 'pending' COMMENT 'pending|paid|failed|refunded',
  `transaction_id` varchar(100) DEFAULT NULL,
  `payment_time` timestamp NULL DEFAULT NULL,
  `refund_time` timestamp NULL DEFAULT NULL,
  `refund_amount` decimal(10,2) DEFAULT NULL,
  `refund_reason` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_payment_orders_order_no` (`order_no`),
  KEY `idx_payment_orders_user_id` (`user_id`),
  KEY `idx_payment_orders_status` (`payment_status`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`package_id`) REFERENCES `premium_packages` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建服务使用日志表
CREATE TABLE IF NOT EXISTS `service_usage_logs` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `service_type` varchar(50) NOT NULL COMMENT 'photo_restore|custom_template|premium_service',
  `service_data` json,
  `usage_count` int DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_service_usage_logs_user_id` (`user_id`),
  KEY `idx_service_usage_logs_type` (`service_type`),
  KEY `idx_service_usage_logs_created_at` (`created_at`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认高级套餐
INSERT INTO `premium_packages` (`id`, `package_name`, `package_type`, `description`, `features`, `price`, `duration`, `storage_size`, `is_active`, `sort_order`) VALUES
('pkg-basic', '基础版', 'memorial', '适合个人使用的基础纪念馆服务', '["100MB存储空间", "基础主题模板", "标准祭扫功能"]', 0.00, 365, 104857600, 1, 1),
('pkg-premium', '高级版', 'memorial', '提供更多定制化功能和存储空间', '["500MB存储空间", "高级主题模板", "定制墓碑样式", "老照片修复", "优先客服支持"]', 99.00, 365, 524288000, 1, 2),
('pkg-vip', '尊享版', 'memorial', '最完整的纪念馆服务体验', '["2GB存储空间", "所有主题模板", "完全定制化", "老照片修复", "专属追思会", "数据备份服务", "专属客服"]', 299.00, 365, 2147483648, 1, 3),
('pkg-storage', '扩展存储包', 'storage', '额外增加1GB存储空间', '["1GB额外存储空间", "永久有效"]', 49.00, 36500, 1073741824, 1, 4);


-- 创建专属服务表
CREATE TABLE IF NOT EXISTS `exclusive_services` (
  `id` varchar(36) NOT NULL,
  `service_name` varchar(100) NOT NULL,
  `service_type` varchar(50) NOT NULL COMMENT 'memorial_service|data_backup|photo_restore|custom_design',
  `description` text,
  `base_price` decimal(10,2) NOT NULL,
  `price_unit` varchar(20) DEFAULT NULL COMMENT 'per_hour|per_service|per_gb',
  `features` json,
  `require_booking` tinyint(1) DEFAULT '0',
  `max_duration` int DEFAULT NULL COMMENT '最大时长（分钟）',
  `is_active` tinyint(1) DEFAULT '1',
  `sort_order` int DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_exclusive_services_type` (`service_type`),
  KEY `idx_exclusive_services_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建服务预订表
CREATE TABLE IF NOT EXISTS `service_bookings` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `service_id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) DEFAULT NULL,
  `booking_type` varchar(50) NOT NULL COMMENT 'memorial_service|data_backup|custom_design',
  `scheduled_time` timestamp NULL DEFAULT NULL,
  `duration` int DEFAULT NULL COMMENT '时长（分钟）',
  `status` varchar(20) DEFAULT 'pending' COMMENT 'pending|confirmed|in_progress|completed|cancelled',
  `total_price` decimal(10,2) DEFAULT NULL,
  `requirements` text,
  `service_data` json,
  `staff_id` varchar(36) DEFAULT NULL,
  `completed_at` timestamp NULL DEFAULT NULL,
  `cancelled_at` timestamp NULL DEFAULT NULL,
  `cancellation_reason` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_service_bookings_user_id` (`user_id`),
  KEY `idx_service_bookings_memorial_id` (`memorial_id`),
  KEY `idx_service_bookings_status` (`status`),
  KEY `idx_service_bookings_scheduled_time` (`scheduled_time`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`service_id`) REFERENCES `exclusive_services` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建数据导出请求表
CREATE TABLE IF NOT EXISTS `data_export_requests` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `export_type` varchar(20) NOT NULL COMMENT 'full|memorial|family',
  `target_id` varchar(36) DEFAULT NULL COMMENT '纪念馆ID或家族ID',
  `export_format` varchar(10) DEFAULT 'zip' COMMENT 'zip|pdf|json',
  `include_media` tinyint(1) DEFAULT '1',
  `encrypted` tinyint(1) DEFAULT '0',
  `encryption_key` varchar(100) DEFAULT NULL,
  `status` varchar(20) DEFAULT 'pending' COMMENT 'pending|processing|completed|failed',
  `file_size` bigint DEFAULT NULL,
  `file_path` varchar(500) DEFAULT NULL,
  `download_url` varchar(500) DEFAULT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  `error_message` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `completed_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_data_export_requests_user_id` (`user_id`),
  KEY `idx_data_export_requests_status` (`status`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建老照片修复请求表
CREATE TABLE IF NOT EXISTS `photo_restore_requests` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) DEFAULT NULL,
  `original_photo_url` varchar(500) NOT NULL,
  `restored_photo_url` varchar(500) DEFAULT NULL,
  `restore_type` varchar(20) NOT NULL COMMENT 'colorize|enhance|repair|all',
  `status` varchar(20) DEFAULT 'pending' COMMENT 'pending|processing|completed|failed',
  `processing_time` int DEFAULT NULL COMMENT '处理时长（秒）',
  `error_message` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `completed_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_photo_restore_requests_user_id` (`user_id`),
  KEY `idx_photo_restore_requests_memorial_id` (`memorial_id`),
  KEY `idx_photo_restore_requests_status` (`status`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建定制设计请求表
CREATE TABLE IF NOT EXISTS `custom_design_requests` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) DEFAULT NULL,
  `design_type` varchar(50) NOT NULL COMMENT 'theme|tombstone|layout|complete',
  `requirements` text,
  `reference_images` json COMMENT '参考图片URL列表',
  `budget` decimal(10,2) DEFAULT NULL,
  `status` varchar(20) DEFAULT 'pending' COMMENT 'pending|in_design|review|completed|cancelled',
  `designer_id` varchar(36) DEFAULT NULL,
  `design_files` json,
  `feedback_count` int DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `completed_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_custom_design_requests_user_id` (`user_id`),
  KEY `idx_custom_design_requests_memorial_id` (`memorial_id`),
  KEY `idx_custom_design_requests_status` (`status`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建服务评价表
CREATE TABLE IF NOT EXISTS `service_reviews` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `booking_id` varchar(36) NOT NULL,
  `rating` int NOT NULL COMMENT '1-5星',
  `comment` text,
  `tags` json COMMENT '评价标签',
  `is_anonymous` tinyint(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_service_reviews_booking_id` (`booking_id`),
  KEY `idx_service_reviews_user_id` (`user_id`),
  KEY `idx_service_reviews_rating` (`rating`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`booking_id`) REFERENCES `service_bookings` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建服务人员表
CREATE TABLE IF NOT EXISTS `service_staff` (
  `id` varchar(36) NOT NULL,
  `name` varchar(50) NOT NULL,
  `role` varchar(50) NOT NULL COMMENT 'designer|coordinator|technician',
  `specialties` json,
  `avatar_url` varchar(255) DEFAULT NULL,
  `bio` text,
  `rating` decimal(3,2) DEFAULT '5.00',
  `review_count` int DEFAULT '0',
  `is_available` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_service_staff_role` (`role`),
  KEY `idx_service_staff_available` (`is_available`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认专属服务
INSERT INTO `exclusive_services` (`id`, `service_name`, `service_type`, `description`, `base_price`, `price_unit`, `features`, `require_booking`, `max_duration`, `is_active`, `sort_order`) VALUES
('svc-memorial', '专属追思会策划', 'memorial_service', '专业团队为您策划和主持线上追思会，提供全程技术支持', 299.00, 'per_service', '["专业主持人", "定制流程", "技术支持", "录制回放", "最多50人参与"]', 1, 180, 1, 1),
('svc-backup', '数据备份导出服务', 'data_backup', '将您的纪念馆数据导出为加密文件，永久保存', 49.00, 'per_service', '["完整数据导出", "加密保护", "多种格式", "7天有效期"]', 0, NULL, 1, 2),
('svc-restore', '老照片AI修复', 'photo_restore', '使用AI技术修复老照片，包括上色、增强、修补等', 19.00, 'per_service', '["AI智能修复", "色彩还原", "清晰度增强", "划痕修复"]', 0, NULL, 1, 3),
('svc-design', '定制设计服务', 'custom_design', '专业设计师为您定制独一无二的纪念馆主题和墓碑', 599.00, 'per_service', '["专业设计师", "一对一沟通", "3次修改机会", "源文件交付"]', 1, NULL, 1, 4);
