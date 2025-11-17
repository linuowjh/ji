-- 创建家族活动表
CREATE TABLE IF NOT EXISTS `family_activities` (
  `id` varchar(36) NOT NULL,
  `family_id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `memorial_id` varchar(36) DEFAULT NULL,
  `activity_type` varchar(30) NOT NULL COMMENT 'worship|join|create_memorial|create_story',
  `content` json DEFAULT NULL,
  `timestamp` timestamp NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_family_activities_family_id` (`family_id`),
  KEY `idx_family_activities_user_id` (`user_id`),
  KEY `idx_family_activities_memorial_id` (`memorial_id`),
  KEY `idx_family_activities_timestamp` (`timestamp`),
  KEY `idx_family_activities_deleted_at` (`deleted_at`),
  FOREIGN KEY (`family_id`) REFERENCES `families` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`memorial_id`) REFERENCES `memorials` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
