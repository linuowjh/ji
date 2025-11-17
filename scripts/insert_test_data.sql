-- 清空现有数据
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE users;
TRUNCATE TABLE memorials;
TRUNCATE TABLE families;
TRUNCATE TABLE family_members;
TRUNCATE TABLE memorial_families;
TRUNCATE TABLE worship_records;
TRUNCATE TABLE prayers;
TRUNCATE TABLE messages;
TRUNCATE TABLE memorial_reminders;
TRUNCATE TABLE family_activities;
SET FOREIGN_KEY_CHECKS = 1;

-- 插入测试用户
INSERT INTO `users` (`id`, `wechat_open_id`, `nickname`, `avatar_url`, `status`, `created_at`, `updated_at`) VALUES
('test-user-1', 'test_openid_1', '张三', 'https://thirdwx.qlogo.cn/mmopen/vi_32/POgEwh4mIHO4nibH0KlMECNjjGxQUq24ZEaGT4poC6icRiccVGKSyXwibcPq4BWmiaIGuG1icwxaQX6grC9VemZoJ8rg/132', 1, NOW(), NOW()),
('test-user-2', 'test_openid_2', '李四', 'https://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoj0hHXhgJNOTSOFsS4uZs8x1ConecaVOB8eIl115xmJZcT4oCicvia7wMEufibKtTLqiaJeanU2Lpg3w/132', 1, NOW(), NOW()),
('test-user-3', 'test_openid_3', '王五', 'https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLL1byctY955Htv9ztzXP4v9qYQRqAzPNLTXzzKvOBz3R6JQ0VfJKJvXkvxNjX8NnKpBMqiaymMGMA/132', 1, NOW(), NOW()),
('test-user-4', 'test_openid_4', '赵六', 'https://thirdwx.qlogo.cn/mmopen/vi_32/ajNVdqHZLLBWribSGCGqRbTq8kKphgzhKqv9rK0xt8B28lGYzo7NE0PlHg5ia1xvhria0AaqGhoXczS2FviaW2dbuw/132', 1, NOW(), NOW());

-- 插入测试纪念馆
INSERT INTO `memorials` (`id`, `creator_id`, `deceased_name`, `birth_date`, `death_date`, `biography`, `avatar_url`, `theme_style`, `tombstone_style`, `privacy_level`, `status`, `created_at`, `updated_at`) VALUES
('test-memorial-1', 'test-user-1', '张老爷子', '1950-01-01', '2020-12-31', '张老爷子是一位慈祥的长者，一生勤劳善良，深受家人和邻里的爱戴。他热爱生活，关爱家人，是我们永远的榜样。', 'https://example.com/avatar1.jpg', 'traditional', 'marble', 1, 1, NOW(), NOW()),
('test-memorial-2', 'test-user-2', '李奶奶', '1955-03-15', '2021-06-20', '李奶奶是一位温柔的母亲，她用自己的爱温暖着整个家庭。她的笑容永远留在我们心中。', 'https://example.com/avatar2.jpg', 'elegant', 'granite', 1, 1, NOW(), NOW()),
('test-memorial-3', 'test-user-3', '王爷爷', '1948-05-20', '2022-03-10', '王爷爷是一位退伍军人，为国家奉献了一生。他的精神永远激励着我们。', 'https://example.com/avatar3.jpg', 'traditional', 'marble', 1, 1, NOW(), NOW());

-- 插入测试家族
INSERT INTO `families` (`id`, `name`, `creator_id`, `description`, `invite_code`, `created_at`, `updated_at`) VALUES
('test-family-1', '张氏家族', 'test-user-1', '张氏家族纪念圈，传承家族情感，共同缅怀先人。', 'ZHANG001', NOW(), NOW()),
('test-family-2', '李氏家族', 'test-user-2', '李氏家族纪念圈，让爱跨越时空。', 'LI002', NOW(), NOW());

-- 插入家族成员
INSERT INTO `family_members` (`id`, `family_id`, `user_id`, `role`, `joined_at`) VALUES
('test-member-1', 'test-family-1', 'test-user-1', 'admin', NOW()),
('test-member-2', 'test-family-1', 'test-user-3', 'member', NOW()),
('test-member-3', 'test-family-1', 'test-user-4', 'member', NOW()),
('test-member-4', 'test-family-2', 'test-user-2', 'admin', NOW()),
('test-member-5', 'test-family-2', 'test-user-3', 'member', NOW());

-- 插入纪念馆家族关联
INSERT INTO `memorial_families` (`id`, `memorial_id`, `family_id`, `created_at`) VALUES
('test-mf-1', 'test-memorial-1', 'test-family-1', NOW()),
('test-mf-2', 'test-memorial-2', 'test-family-2', NOW()),
('test-mf-3', 'test-memorial-3', 'test-family-1', NOW());

-- 插入祭扫记录
INSERT INTO `worship_records` (`id`, `memorial_id`, `user_id`, `worship_type`, `content`, `created_at`, `updated_at`) VALUES
('test-worship-1', 'test-memorial-1', 'test-user-1', 'flower', '{"flower_type": "菊花", "count": 3, "message": "爷爷，我们想您了"}', NOW(), NOW()),
('test-worship-2', 'test-memorial-1', 'test-user-3', 'candle', '{"candle_type": "红烛", "duration": 24, "message": "为爷爷点亮心灯"}', DATE_SUB(NOW(), INTERVAL 1 HOUR), DATE_SUB(NOW(), INTERVAL 1 HOUR)),
('test-worship-3', 'test-memorial-2', 'test-user-2', 'incense', '{"incense_count": 3, "message": "妈妈，愿您在天堂安好"}', DATE_SUB(NOW(), INTERVAL 2 HOUR), DATE_SUB(NOW(), INTERVAL 2 HOUR)),
('test-worship-4', 'test-memorial-3', 'test-user-3', 'flower', '{"flower_type": "白菊", "count": 5, "message": "王爷爷，永远怀念您"}', DATE_SUB(NOW(), INTERVAL 3 HOUR), DATE_SUB(NOW(), INTERVAL 3 HOUR));

-- 插入祈福
INSERT INTO `prayers` (`id`, `memorial_id`, `user_id`, `content`, `is_public`, `created_at`, `updated_at`) VALUES
('test-prayer-1', 'test-memorial-1', 'test-user-1', '爷爷，愿您在天堂安好，我们会好好生活，不让您担心。', 1, NOW(), NOW()),
('test-prayer-2', 'test-memorial-2', 'test-user-2', '妈妈，谢谢您给我们的爱，我们永远爱您。', 1, NOW(), NOW());

-- 插入留言
INSERT INTO `messages` (`id`, `memorial_id`, `user_id`, `message_type`, `content`, `created_at`, `updated_at`) VALUES
('test-message-1', 'test-memorial-1', 'test-user-1', 'text', '爷爷，今天是您的生日，我们全家都在想念您。您教给我们的做人道理，我们会一直记在心里。', NOW(), NOW()),
('test-message-2', 'test-memorial-2', 'test-user-2', 'text', '妈妈，今天路过您最喜欢的花园，看到满园的花开，就想起了您的笑容。', NOW(), NOW());

-- 插入纪念日提醒
INSERT INTO `memorial_reminders` (`id`, `memorial_id`, `reminder_type`, `reminder_date`, `title`, `content`, `is_active`, `created_at`, `updated_at`) VALUES
('test-reminder-1', 'test-memorial-1', 'birthday', '2025-01-01', '张老爷子生日', '今天是张老爷子的生日，让我们一起为他献花祈福。', 1, NOW(), NOW()),
('test-reminder-2', 'test-memorial-1', 'death_anniversary', '2025-12-31', '张老爷子忌日', '今天是张老爷子的忌日，愿他在天堂安好。', 1, NOW(), NOW()),
('test-reminder-3', 'test-memorial-2', 'birthday', '2025-03-15', '李奶奶生日', '今天是李奶奶的生日，让我们一起缅怀她的慈爱。', 1, NOW(), NOW());

-- 插入家族活动
INSERT INTO `family_activities` (`id`, `family_id`, `user_id`, `memorial_id`, `activity_type`, `content`, `timestamp`, `created_at`) VALUES
('test-activity-1', 'test-family-1', 'test-user-1', 'test-memorial-1', 'worship', '{"worship_type": "flower", "message": "献上鲜花"}', NOW(), NOW()),
('test-activity-2', 'test-family-1', 'test-user-3', NULL, 'join', '{"message": "加入了家族圈"}', DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY)),
('test-activity-3', 'test-family-1', 'test-user-1', 'test-memorial-3', 'create_memorial', '{"memorial_name": "王爷爷"}', DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY)),
('test-activity-4', 'test-family-2', 'test-user-2', 'test-memorial-2', 'worship', '{"worship_type": "incense", "message": "敬献香火"}', DATE_SUB(NOW(), INTERVAL 1 HOUR), DATE_SUB(NOW(), INTERVAL 1 HOUR));
