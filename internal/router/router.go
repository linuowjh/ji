package router

import (
	"yun-nian-memorial/internal/config"
	"yun-nian-memorial/internal/controllers"
	"yun-nian-memorial/internal/middleware"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 基础中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.XSSProtection())
	r.Use(middleware.SQLInjectionProtection())
	r.Use(middleware.SecurityLogger())
	r.Use(middleware.UserAgentFilter())
	r.Use(middleware.RequestSizeLimit(cfg.Security.MaxRequestSize))

	// 初始化服务
	userService := services.NewUserService(db, cfg)
	memorialService := services.NewMemorialService(db)
	mediaService := services.NewMediaService(db, "uploads") // 上传目录
	worshipService := services.NewWorshipService(db)
	albumService := services.NewAlbumService(db)
	lifeStoryService := services.NewLifeStoryService(db)
	memorialServiceService := services.NewMemorialServiceService(db)
	familyService := services.NewFamilyService(db)
	privacyService := services.NewPrivacyService(db)
	adminService := services.NewAdminService(db)

	// 设置服务依赖关系（避免循环依赖）
	worshipService.SetFamilyService(familyService)

	// 初始化控制器
	userController := controllers.NewUserController(userService)
	memorialController := controllers.NewMemorialController(memorialService)
	mediaController := controllers.NewMediaController(mediaService)
	worshipController := controllers.NewWorshipController(worshipService)
	albumController := controllers.NewAlbumController(albumService)
	lifeStoryController := controllers.NewLifeStoryController(lifeStoryService)
	memorialServiceController := controllers.NewMemorialServiceController(memorialServiceService)
	familyController := controllers.NewFamilyController(familyService)
	privacyController := controllers.NewPrivacyController(privacyService)
	adminController := controllers.NewAdminController(adminService)

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "云念纪念馆服务运行正常",
			"version": "1.0.0",
		})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由（无需登录）
		auth := api.Group("/auth")
		{
			auth.POST("/wechat-login", userController.WechatLogin)
		}

		// 需要认证的路由
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			// 用户相关路由
			users := protected.Group("/users")
			{
				users.GET("/profile", userController.GetUserInfo)
				users.PUT("/profile", userController.UpdateUserInfo)
				users.GET("/memorials", userController.GetUserMemorials)
				users.GET("/worship-records", userController.GetUserWorshipRecords)

				// 个人中心功能
				users.GET("/dashboard", userController.GetUserDashboard)
				users.GET("/statistics", userController.GetUserStatistics)
				users.GET("/activities", userController.GetUserRecentActivities)
				users.GET("/memorial-details", userController.GetUserMemorialDetails)
				users.GET("/families", userController.GetUserFamilies)
				users.GET("/memorials/:memorial_id/visitors", userController.GetMemorialVisitors)
				users.GET("/reminders/upcoming", userController.GetUpcomingReminders)
			}

			// 纪念馆相关路由
			memorials := protected.Group("/memorials")
			{
				memorials.GET("/recent", memorialController.GetRecentMemorials) // 最近访问的纪念馆
				memorials.GET("/", memorialController.GetMemorialList)
				memorials.POST("/", memorialController.CreateMemorial)
				memorials.GET("/:id", memorialController.GetMemorial)
				memorials.PUT("/:id", memorialController.UpdateMemorial)
				memorials.DELETE("/:id", memorialController.DeleteMemorial)
				memorials.GET("/:id/visitors", memorialController.GetMemorialVisitors)

				// 墓碑定制相关路由
				memorials.PUT("/:id/tombstone-style", memorialController.UpdateTombstoneStyle)
				memorials.PUT("/:id/epitaph", memorialController.UpdateEpitaph)
			}

			// 样式和工具相关路由
			styles := protected.Group("/styles")
			{
				styles.GET("/tombstones", memorialController.GetTombstoneStyles)
				styles.GET("/themes", memorialController.GetThemeStyles)
			}

			// 工具相关路由
			tools := protected.Group("/tools")
			{
				tools.POST("/calligraphy", memorialController.GenerateCalligraphy)
				tools.POST("/handwriting", memorialController.ProcessHandwriting)
			}

			// 媒体文件相关路由
			media := protected.Group("/media")
			{
				media.POST("/upload", mediaController.Upload) // 通用上传接口
				media.POST("/upload/image", mediaController.UploadImage)
				media.POST("/upload/video", mediaController.UploadVideo)
				media.POST("/upload/audio", mediaController.UploadAudio)
				media.GET("/memorials/:memorial_id/files", mediaController.GetMediaFiles)
				media.GET("/memorials/:memorial_id/stats", mediaController.GetMediaFileStats)
				media.PUT("/files/:id", mediaController.UpdateMediaFile)
				media.DELETE("/files/:id", mediaController.DeleteMediaFile)
			}

			// 祭扫相关路由
			worship := protected.Group("/worship")
			{
				// 传统祭扫功能
				worship.POST("/memorials/:memorial_id/flowers", worshipController.OfferFlowers)
				worship.POST("/memorials/:memorial_id/candles", worshipController.LightCandle)
				worship.PUT("/memorials/:memorial_id/candles/renew", worshipController.RenewCandle)
				worship.GET("/memorials/:memorial_id/candles/status", worshipController.GetCandleStatus)
				worship.POST("/memorials/:memorial_id/incense", worshipController.OfferIncense)
				worship.POST("/memorials/:memorial_id/tributes", worshipController.OfferTribute)

				// 祈福和留言功能
				worship.POST("/memorials/:memorial_id/prayers", worshipController.CreatePrayer)
				worship.POST("/memorials/:memorial_id/messages", worshipController.CreateMessage)
				worship.POST("/scheduled-prayers", worshipController.CreateScheduledPrayer)

				// 祈福卡功能
				worship.GET("/prayer-card-templates", worshipController.GetPrayerCardTemplates)
				worship.POST("/generate-prayer-card", worshipController.GeneratePrayerCard)
				worship.GET("/popular-prayer-contents", worshipController.GetPopularPrayerContents)

				// 智能功能
				worship.POST("/analyze-emotion", worshipController.AnalyzeMessageEmotion)
				worship.GET("/reply-suggestions", worshipController.GetMessageReplySuggestions)
				worship.GET("/memorials/:memorial_id/message-tips", worshipController.GetMessageCreationTips)
				worship.GET("/memorials/:memorial_id/message-analytics", worshipController.GetMemorialMessageAnalytics)

				// 内容审核
				worship.POST("/messages/:message_id/moderate", worshipController.ModerateMessage)

				// 查询功能
				worship.GET("/memorials/:memorial_id/records", worshipController.GetWorshipRecords)
				worship.GET("/memorials/:memorial_id/prayer-wall", worshipController.GetPrayerWall)
				worship.GET("/memorials/:memorial_id/time-messages", worshipController.GetTimeMessages)
				worship.GET("/memorials/:memorial_id/statistics", worshipController.GetWorshipStatistics)
				worship.GET("/memorials/:memorial_id/detailed-statistics", worshipController.GetDetailedWorshipStatistics)
				worship.GET("/memorials/:memorial_id/report", worshipController.GenerateWorshipReport)
				worship.GET("/user/history", worshipController.GetUserWorshipHistory)
				worship.GET("/user/behavior-analysis", worshipController.AnalyzeUserWorshipBehavior)
			}

			// 纪念相册相关路由
			albums := protected.Group("/albums")
			{
				albums.POST("/memorials/:memorial_id", albumController.CreateAlbum)
				albums.GET("/memorials/:memorial_id", albumController.GetAlbums)
				albums.GET("/:album_id", albumController.GetAlbum)
				albums.PUT("/:album_id", albumController.UpdateAlbum)
				albums.DELETE("/:album_id", albumController.DeleteAlbum)

				// 相册照片管理
				albums.POST("/:album_id/photos", albumController.AddPhoto)
				albums.PUT("/photos/:photo_id", albumController.UpdatePhoto)
				albums.DELETE("/photos/:photo_id", albumController.DeletePhoto)
			}

			// 生平故事相关路由
			stories := protected.Group("/stories")
			{
				stories.POST("/memorials/:memorial_id", lifeStoryController.CreateLifeStory)
				stories.GET("/memorials/:memorial_id", lifeStoryController.GetLifeStories)
				stories.GET("/memorials/:memorial_id/by-category", lifeStoryController.GetStoriesByCategory)
				stories.GET("/:story_id", lifeStoryController.GetLifeStory)
				stories.PUT("/:story_id", lifeStoryController.UpdateLifeStory)
				stories.DELETE("/:story_id", lifeStoryController.DeleteLifeStory)
			}

			// 时间轴相关路由
			timelines := protected.Group("/timelines")
			{
				timelines.POST("/memorials/:memorial_id", lifeStoryController.CreateTimeline)
				timelines.GET("/memorials/:memorial_id", lifeStoryController.GetTimeline)
				timelines.DELETE("/:timeline_id", lifeStoryController.DeleteTimeline)
			}

			// 线上追思会相关路由
			memorialServices := protected.Group("/memorial-services")
			{
				// 追思会管理
				memorialServices.POST("/memorials/:memorial_id", memorialServiceController.CreateMemorialService)
				memorialServices.GET("/memorials/:memorial_id", memorialServiceController.GetMemorialServices)
				memorialServices.GET("/:service_id", memorialServiceController.GetMemorialService)
				memorialServices.PUT("/:service_id", memorialServiceController.UpdateMemorialService)
				memorialServices.DELETE("/:service_id", memorialServiceController.DeleteMemorialService)

				// 追思会控制
				memorialServices.POST("/:service_id/start", memorialServiceController.StartService)
				memorialServices.POST("/:service_id/end", memorialServiceController.EndService)
				memorialServices.POST("/:service_id/join", memorialServiceController.JoinService)
				memorialServices.POST("/:service_id/leave", memorialServiceController.LeaveService)

				// 参与者管理
				memorialServices.POST("/:service_id/invite", memorialServiceController.InviteParticipants)
				memorialServices.POST("/invitations/:invitation_id/respond", memorialServiceController.RespondToInvitation)

				// 聊天功能
				memorialServices.POST("/:service_id/chat", memorialServiceController.SendChatMessage)
				memorialServices.GET("/:service_id/chat", memorialServiceController.GetChatMessages)
			}

			// 家族相关路由
			families := protected.Group("/families")
			{
				// 家族圈管理
				families.GET("/", familyController.GetFamilies)
				families.POST("/", familyController.CreateFamily)
				families.GET("/:family_id", familyController.GetFamily)
				families.PUT("/:family_id", familyController.UpdateFamily)
				families.DELETE("/:family_id", familyController.DeleteFamily)

				// 成员管理
				families.GET("/:family_id/members", familyController.GetFamilyMembers)
				families.POST("/:family_id/invite", familyController.InviteMembers)
				families.DELETE("/:family_id/members/:member_id", familyController.RemoveMember)
				families.PUT("/:family_id/members/:member_id/role", familyController.SetMemberRole)
				families.POST("/:family_id/leave", familyController.LeaveFamily)

				// 邀请管理
				families.POST("/join-by-code", familyController.JoinFamilyByCode)
				families.POST("/invitations/:invitation_id/respond", familyController.RespondToInvitation)

				// 家族活动
				families.GET("/:family_id/activities", familyController.GetFamilyActivities)

				// 纪念馆关联
				families.POST("/:family_id/memorials", familyController.AddMemorialToFamily)
				families.DELETE("/:family_id/memorials/:memorial_id", familyController.RemoveMemorialFromFamily)

				// 纪念日提醒
				families.POST("/:family_id/reminders", familyController.SetMemorialReminder)
				families.GET("/:family_id/reminders", familyController.GetFamilyReminders)
				families.GET("/:family_id/reminders/upcoming", familyController.GetUpcomingReminders)
				families.DELETE("/:family_id/reminders/:reminder_id", familyController.DeleteReminder)

				// 集体祭扫
				families.POST("/:family_id/collective-worship", familyController.InitiateCollectiveWorship)
				families.POST("/:family_id/collective-worship/:activity_id/join", familyController.JoinCollectiveWorship)

				// 家族谱系
				families.POST("/:family_id/genealogy", familyController.CreateGenealogy)
				families.GET("/:family_id/genealogy", familyController.GetFamilyGenealogy)
				families.PUT("/:family_id/genealogy/:genealogy_id", familyController.UpdateGenealogy)
				families.DELETE("/:family_id/genealogy/:genealogy_id", familyController.DeleteGenealogy)

				// 家族故事
				families.POST("/:family_id/stories", familyController.CreateFamilyStory)
				families.GET("/:family_id/stories", familyController.GetFamilyStories)
				families.GET("/:family_id/stories/:story_id", familyController.GetFamilyStory)
				families.PUT("/:family_id/stories/:story_id", familyController.UpdateFamilyStory)
				families.DELETE("/:family_id/stories/:story_id", familyController.DeleteFamilyStory)

				// 家族传统
				families.POST("/:family_id/traditions", familyController.CreateFamilyTradition)
				families.GET("/:family_id/traditions", familyController.GetFamilyTraditions)
				families.PUT("/:family_id/traditions/:tradition_id", familyController.UpdateFamilyTradition)
				families.DELETE("/:family_id/traditions/:tradition_id", familyController.DeleteFamilyTradition)
			}

			// 隐私设置相关路由
			privacy := protected.Group("/privacy")
			{
				// 纪念馆隐私设置
				privacy.POST("/memorials/settings", privacyController.SetMemorialPrivacy)
				privacy.GET("/memorials/:memorial_id/settings", privacyController.GetMemorialPrivacySettings)

				// 访问权限检查
				privacy.GET("/memorials/:memorial_id/access", privacyController.CheckUserAccess)

				// 访问申请
				privacy.POST("/memorials/:memorial_id/request-access", privacyController.RequestAccess)
				privacy.GET("/memorials/:memorial_id/access-requests", privacyController.GetAccessRequests)
				privacy.POST("/access-requests/:request_id/handle", privacyController.HandleAccessRequest)

				// 黑名单管理
				privacy.POST("/memorials/:memorial_id/blacklist/:user_id", privacyController.AddToBlacklist)
				privacy.DELETE("/memorials/:memorial_id/blacklist/:user_id", privacyController.RemoveFromBlacklist)
			}

			// 系统管理相关路由（需要管理员权限）
			admin := protected.Group("/admin")
			// 为管理员接口添加额外的安全保护
			admin.Use(middleware.IPWhitelist(cfg.Security.AdminIPWhitelist))
			{
				// 用户管理
				admin.GET("/users", adminController.GetUserList)
				admin.GET("/users/:user_id", adminController.GetUserDetail)
				admin.POST("/users/manage", adminController.ManageUser)

				// 内容审核
				admin.GET("/content/pending", adminController.GetPendingContent)
				admin.POST("/content/moderate", adminController.ModerateContent)
				admin.POST("/content/batch-moderate", adminController.BatchModerateContent)

				// 系统统计
				admin.GET("/stats", adminController.GetSystemStats)
			}
		}
	}

	return r
}
