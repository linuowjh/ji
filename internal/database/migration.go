package database

import (
	"fmt"
	"log"
	"strings"
	"yun-nian-memorial/internal/models"

	"gorm.io/gorm"
)

// Migration 数据库迁移管理器
type Migration struct {
	db *gorm.DB
}

// NewMigration 创建迁移管理器
func NewMigration(db *gorm.DB) *Migration {
	return &Migration{db: db}
}

// AutoMigrate 自动迁移所有表
func (m *Migration) AutoMigrate() error {
	log.Println("开始执行数据库自动迁移...")

	// 定义所有需要迁移的模型
	models := []interface{}{
		&models.User{},
		&models.Memorial{},
		&models.WorshipRecord{},
		&models.Family{},
		&models.FamilyMember{},
		&models.MediaFile{},
		&models.Prayer{},
		&models.Message{},
		&models.MemorialReminder{},
		&models.VisitorRecord{},
		&models.MemorialFamily{},
		&models.Album{},
		&models.AlbumPhoto{},
		&models.LifeStory{},
		&models.LifeStoryMedia{},
		&models.Timeline{},
		&models.MemorialService{},
		&models.MemorialServiceParticipant{},
		&models.ServiceActivity{},
		&models.ServiceInvitation{},
		&models.ServiceRecording{},
		&models.ServiceChat{},
		&models.FamilyGenealogy{},
		&models.FamilyStory{},
		&models.FamilyTradition{},
		&models.VisitorPermissionSetting{},
		&models.VisitorBlacklist{},
		&models.AccessRequest{},
		// 系统配置和维护相关模型
		&models.SystemConfig{},
		&models.FestivalConfig{},
		&models.TemplateConfig{},
		&models.DataBackup{},
		&models.SystemLog{},
		&models.SystemMonitor{},
		// 增值服务相关模型
		&models.PremiumPackage{},
		&models.UserSubscription{},
		&models.MemorialUpgrade{},
		&models.CustomTemplate{},
		&models.StorageUsage{},
		&models.PaymentOrder{},
		&models.ServiceUsageLog{},
		// 专属服务相关模型
		&models.ExclusiveService{},
		&models.ServiceBooking{},
		&models.DataExportRequest{},
		&models.PhotoRestoreRequest{},
		&models.CustomDesignRequest{},
		&models.ServiceReview{},
		&models.ServiceStaff{},
	}

	// 执行自动迁移
	for _, model := range models {
		if err := m.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("迁移模型 %T 失败: %v", model, err)
		}
		log.Printf("成功迁移模型: %T", model)
	}

	log.Println("数据库自动迁移完成")
	return nil
}

// CreateIndexes 创建额外的索引
func (m *Migration) CreateIndexes() error {
	log.Println("开始创建数据库索引...")

	// MySQL兼容的索引定义（不使用IF NOT EXISTS）
	indexes := []string{
		// 用户表索引
		"CREATE INDEX idx_users_status ON users(status)",
		"CREATE INDEX idx_users_created_at ON users(created_at)",

		// 纪念馆表索引
		"CREATE INDEX idx_memorials_creator_status ON memorials(creator_id, status)",
		"CREATE INDEX idx_memorials_privacy_status ON memorials(privacy_level, status)",
		"CREATE INDEX idx_memorials_created_at ON memorials(created_at)",

		// 祭扫记录表索引
		"CREATE INDEX idx_worship_memorial_time ON worship_records(memorial_id, created_at)",
		"CREATE INDEX idx_worship_user_time ON worship_records(user_id, created_at)",
		"CREATE INDEX idx_worship_type ON worship_records(worship_type)",

		// 家族表索引
		"CREATE INDEX idx_families_created_at ON families(created_at)",

		// 家族成员表索引
		"CREATE INDEX idx_family_members_family_role ON family_members(family_id, role)",

		// 媒体文件表索引
		"CREATE INDEX idx_media_files_type_memorial ON media_files(file_type, memorial_id)",

		// 祈福表索引
		"CREATE INDEX idx_prayers_memorial_public ON prayers(memorial_id, is_public)",
		"CREATE INDEX idx_prayers_created_at ON prayers(created_at)",

		// 留言表索引
		"CREATE INDEX idx_messages_memorial_time ON messages(memorial_id, created_at)",
		"CREATE INDEX idx_messages_type ON messages(message_type)",

		// 纪念日提醒表索引
		"CREATE INDEX idx_memorial_reminders_date_active ON memorial_reminders(reminder_date, is_active)",
		"CREATE INDEX idx_memorial_reminders_type ON memorial_reminders(reminder_type)",

		// 访客记录表索引
		"CREATE INDEX idx_visitor_records_time ON visitor_records(visit_time)",
		"CREATE INDEX idx_visitor_records_memorial_time ON visitor_records(memorial_id, visit_time)",

		// 相册表索引
		"CREATE INDEX idx_albums_memorial_id ON albums(memorial_id)",
		"CREATE INDEX idx_albums_created_at ON albums(created_at)",

		// 相册照片表索引
		"CREATE INDEX idx_album_photos_album_id ON album_photos(album_id)",
		"CREATE INDEX idx_album_photos_sort_order ON album_photos(album_id, sort_order)",

		// 生平故事表索引
		"CREATE INDEX idx_life_stories_memorial_id ON life_stories(memorial_id)",
		"CREATE INDEX idx_life_stories_category ON life_stories(memorial_id, category)",
		"CREATE INDEX idx_life_stories_date ON life_stories(story_date)",
		"CREATE INDEX idx_life_stories_author ON life_stories(author_id)",

		// 生平故事媒体表索引
		"CREATE INDEX idx_life_story_media_story_id ON life_story_media(life_story_id)",
		"CREATE INDEX idx_life_story_media_sort_order ON life_story_media(life_story_id, sort_order)",

		// 时间轴表索引
		"CREATE INDEX idx_timelines_memorial_id ON timelines(memorial_id)",
		"CREATE INDEX idx_timelines_event_date ON timelines(memorial_id, event_date)",
		"CREATE INDEX idx_timelines_event_type ON timelines(event_type)",

		// 追思会表索引
		"CREATE INDEX idx_memorial_services_memorial_id ON memorial_services(memorial_id)",
		"CREATE INDEX idx_memorial_services_host_id ON memorial_services(host_id)",
		"CREATE INDEX idx_memorial_services_start_time ON memorial_services(start_time)",
		"CREATE INDEX idx_memorial_services_status ON memorial_services(status)",

		// 追思会参与者表索引
		"CREATE INDEX idx_service_participants_service_id ON memorial_service_participants(service_id)",
		"CREATE INDEX idx_service_participants_user_id ON memorial_service_participants(user_id)",
		"CREATE INDEX idx_service_participants_role ON memorial_service_participants(service_id, role)",

		// 追思会活动表索引
		"CREATE INDEX idx_service_activities_service_id ON service_activities(service_id)",
		"CREATE INDEX idx_service_activities_timestamp ON service_activities(service_id, timestamp)",
		"CREATE INDEX idx_service_activities_type ON service_activities(activity_type)",

		// 追思会邀请表索引
		"CREATE INDEX idx_service_invitations_service_id ON service_invitations(service_id)",
		"CREATE INDEX idx_service_invitations_invitee_id ON service_invitations(invitee_id)",
		"CREATE INDEX idx_service_invitations_status ON service_invitations(status)",

		// 追思会录制表索引
		"CREATE INDEX idx_service_recordings_service_id ON service_recordings(service_id)",
		"CREATE INDEX idx_service_recordings_status ON service_recordings(status)",

		// 追思会聊天表索引
		"CREATE INDEX idx_service_chats_service_id ON service_chats(service_id)",
		"CREATE INDEX idx_service_chats_timestamp ON service_chats(service_id, timestamp)",
	}

	successCount := 0
	skipCount := 0
	
	for _, indexSQL := range indexes {
		if err := m.db.Exec(indexSQL).Error; err != nil {
			// 如果索引已存在，跳过错误
			if strings.Contains(err.Error(), "Duplicate key name") || strings.Contains(err.Error(), "already exists") {
				skipCount++
			} else {
				log.Printf("创建索引失败: %s, 错误: %v", indexSQL, err)
			}
		} else {
			successCount++
		}
	}

	log.Printf("数据库索引创建完成 - 成功: %d, 跳过: %d", successCount, skipCount)
	return nil
}

// SeedData 插入种子数据
func (m *Migration) SeedData() error {
	log.Println("开始插入种子数据...")

	// 检查是否已有数据
	var userCount int64
	m.db.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("数据库已有数据，跳过种子数据插入")
		return nil
	}

	// 创建测试用户
	users := []models.User{
		{
			ID:           "test-user-1",
			WechatOpenID: "test_openid_1",
			Nickname:     "张三",
			AvatarURL:    "https://example.com/avatar1.jpg",
			Status:       1,
		},
		{
			ID:           "test-user-2",
			WechatOpenID: "test_openid_2",
			Nickname:     "李四",
			AvatarURL:    "https://example.com/avatar2.jpg",
			Status:       1,
		},
		{
			ID:           "test-user-3",
			WechatOpenID: "test_openid_3",
			Nickname:     "王五",
			AvatarURL:    "https://example.com/avatar3.jpg",
			Status:       1,
		},
	}

	for _, user := range users {
		if err := m.db.Create(&user).Error; err != nil {
			return fmt.Errorf("创建测试用户失败: %v", err)
		}
	}

	// 创建测试纪念馆
	memorials := []models.Memorial{
		{
			ID:             "test-memorial-1",
			CreatorID:      "test-user-1",
			DeceasedName:   "张老爷子",
			Biography:      "张老爷子是一位慈祥的长者，一生勤劳善良，深受家人和邻里的爱戴。",
			ThemeStyle:     "traditional",
			TombstoneStyle: "marble",
			PrivacyLevel:   1,
			Status:         1,
		},
		{
			ID:             "test-memorial-2",
			CreatorID:      "test-user-2",
			DeceasedName:   "李奶奶",
			Biography:      "李奶奶是一位温柔的母亲，她用自己的爱温暖着整个家庭。",
			ThemeStyle:     "elegant",
			TombstoneStyle: "granite",
			PrivacyLevel:   1,
			Status:         1,
		},
	}

	for _, memorial := range memorials {
		if err := m.db.Create(&memorial).Error; err != nil {
			return fmt.Errorf("创建测试纪念馆失败: %v", err)
		}
	}

	// 创建测试家族
	families := []models.Family{
		{
			ID:          "test-family-1",
			Name:        "张氏家族",
			CreatorID:   "test-user-1",
			Description: "张氏家族纪念圈，传承家族情感，共同缅怀先人。",
			InviteCode:  "ZHANG001",
		},
		{
			ID:          "test-family-2",
			Name:        "李氏家族",
			CreatorID:   "test-user-2",
			Description: "李氏家族纪念圈，让爱跨越时空。",
			InviteCode:  "LI002",
		},
	}

	for _, family := range families {
		if err := m.db.Create(&family).Error; err != nil {
			return fmt.Errorf("创建测试家族失败: %v", err)
		}
	}

	log.Println("种子数据插入完成")
	return nil
}

// DropAllTables 删除所有表（危险操作，仅用于开发环境）
func (m *Migration) DropAllTables() error {
	log.Println("警告：正在删除所有数据表...")

	tables := []string{
		"service_chats",
		"service_recordings",
		"service_invitations",
		"service_activities",
		"memorial_service_participants",
		"memorial_services",
		"timelines",
		"life_story_media",
		"life_stories",
		"album_photos",
		"albums",
		"visitor_records",
		"memorial_families",
		"memorial_reminders",
		"messages",
		"prayers",
		"media_files",
		"family_members",
		"families",
		"worship_records",
		"memorials",
		"users",
	}

	for _, table := range tables {
		if err := m.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)).Error; err != nil {
			log.Printf("删除表 %s 失败: %v", table, err)
		} else {
			log.Printf("成功删除表: %s", table)
		}
	}

	log.Println("所有数据表删除完成")
	return nil
}

// Reset 重置数据库（删除所有表并重新创建）
func (m *Migration) Reset() error {
	log.Println("开始重置数据库...")

	if err := m.DropAllTables(); err != nil {
		return err
	}

	if err := m.AutoMigrate(); err != nil {
		return err
	}

	if err := m.CreateIndexes(); err != nil {
		return err
	}

	if err := m.SeedData(); err != nil {
		return err
	}

	log.Println("数据库重置完成")
	return nil
}