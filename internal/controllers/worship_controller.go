package controllers

import (
	"net/http"
	"strconv"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type WorshipController struct {
	worshipService *services.WorshipService
}

func NewWorshipController(worshipService *services.WorshipService) *WorshipController {
	return &WorshipController{
		worshipService: worshipService,
	}
}

// OfferFlowers 献花
func (c *WorshipController) OfferFlowers(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	var req services.OfferFlowersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.worshipService.OfferFlowers(userID.(string), memorialID, &req)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "献花成功",
	})
}

// LightCandle 点烛
func (c *WorshipController) LightCandle(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	var req services.LightCandleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.worshipService.LightCandle(userID.(string), memorialID, &req)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "点烛成功",
	})
}

// RenewCandle 续烛
func (c *WorshipController) RenewCandle(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	var req struct {
		AdditionalMinutes int `json:"additional_minutes" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.worshipService.RenewCandle(userID.(string), memorialID, req.AdditionalMinutes)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusBadRequest, APIResponse{
				Code:    1001,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "续烛成功",
	})
}

// GetCandleStatus 获取蜡烛状态
func (c *WorshipController) GetCandleStatus(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	status, err := c.worshipService.GetActiveCandleStatus(userID.(string), memorialID)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    status,
	})
}

// OfferIncense 上香
func (c *WorshipController) OfferIncense(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	var req services.OfferIncenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.worshipService.OfferIncense(userID.(string), memorialID, &req)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "上香成功",
	})
}

// OfferTribute 供奉供品
func (c *WorshipController) OfferTribute(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	var req services.OfferTributeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.worshipService.OfferTribute(userID.(string), memorialID, &req)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "供奉成功",
	})
}

// CreatePrayer 创建祈福
func (c *WorshipController) CreatePrayer(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	var req services.CreatePrayerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	prayer, err := c.worshipService.CreatePrayer(userID.(string), memorialID, &req)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "祈福成功",
		Data:    prayer,
	})
}

// CreateMessage 创建留言
func (c *WorshipController) CreateMessage(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	var req services.CreateMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	message, err := c.worshipService.CreateMessage(userID.(string), memorialID, &req)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusBadRequest, APIResponse{
				Code:    1001,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "留言成功",
		Data:    message,
	})
}

// GetWorshipRecords 获取祭扫记录
func (c *WorshipController) GetWorshipRecords(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	records, total, err := c.worshipService.GetWorshipRecords(userID.(string), memorialID, page, pageSize)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	// 转换数据格式，添加友好的显示字段
	formattedRecords := make([]gin.H, 0, len(records))
	worshipTypeMap := map[string]string{
		"flower":  "献了鲜花",
		"candle":  "点燃了蜡烛",
		"incense": "敬献了香火",
		"tribute": "供奉了供品",
		"prayer":  "送上了祈福",
	}

	for _, record := range records {
		formattedRecord := gin.H{
			"id":              record.ID,
			"memorialId":      record.MemorialID,
			"userId":          record.UserID,
			"worshipType":     record.WorshipType,
			"worshipTypeText": worshipTypeMap[record.WorshipType],
			"content":         record.Content,
			"createdAt":       record.CreatedAt,
			"userName":        record.User.Nickname,
			"userAvatar":      record.User.AvatarURL,
		}
		formattedRecords = append(formattedRecords, formattedRecord)
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data: gin.H{
			"list":      formattedRecords,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetPrayerWall 获取祈福墙
func (c *WorshipController) GetPrayerWall(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	prayers, total, err := c.worshipService.GetPrayerWall(userID.(string), memorialID, page, pageSize)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data: gin.H{
			"list":      prayers,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetTimeMessages 获取时光信箱
func (c *WorshipController) GetTimeMessages(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	messages, total, err := c.worshipService.GetTimeMessages(userID.(string), memorialID, page, pageSize)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data: gin.H{
			"list":      messages,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetWorshipStatistics 获取祭扫统计
func (c *WorshipController) GetWorshipStatistics(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	stats, err := c.worshipService.GetWorshipStatistics(userID.(string), memorialID)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    stats,
	})
}

// GetPrayerCardTemplates 获取祈福卡模板
func (c *WorshipController) GetPrayerCardTemplates(ctx *gin.Context) {
	templates := c.worshipService.GetPrayerCardTemplates()

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    templates,
	})
}

// GeneratePrayerCard 生成祈福卡
func (c *WorshipController) GeneratePrayerCard(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	_ = userID // TODO: 可能需要使用 userID

	var req struct {
		TemplateID string `json:"template_id" binding:"required"`
		Content    string `json:"content" binding:"required"`
		UserName   string `json:"user_name"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	cardURL, err := c.worshipService.GeneratePrayerCard(req.TemplateID, req.Content, req.UserName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "生成成功",
		Data: gin.H{
			"card_url": cardURL,
		},
	})
}

// ModerateMessage 审核留言
func (c *WorshipController) ModerateMessage(ctx *gin.Context) {
	messageID := ctx.Param("message_id")
	if messageID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "留言ID不能为空",
		})
		return
	}

	status, err := c.worshipService.ModerateMessage(messageID)
	if err != nil {
		if err.Error() == "留言不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "审核完成",
		Data:    status,
	})
}

// GetUserWorshipHistory 获取用户祭扫历史
func (c *WorshipController) GetUserWorshipHistory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	history, err := c.worshipService.GetUserWorshipHistory(userID.(string), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    history,
	})
}

// CreateScheduledPrayer 创建定时祈福
func (c *WorshipController) CreateScheduledPrayer(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	var req services.ScheduledPrayerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.worshipService.CreateScheduledPrayer(userID.(string), &req)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusBadRequest, APIResponse{
				Code:    1001,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "定时祈福创建成功",
	})
}

// GetPopularPrayerContents 获取热门祈福内容
func (c *WorshipController) GetPopularPrayerContents(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	contents, err := c.worshipService.GetPopularPrayerContents(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data: gin.H{
			"contents": contents,
		},
	})
}

// AnalyzeMessageEmotion 分析留言情感
func (c *WorshipController) AnalyzeMessageEmotion(ctx *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	result, err := c.worshipService.AnalyzeMessageEmotion(req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "分析完成",
		Data:    result,
	})
}

// GetMessageReplySuggestions 获取留言回复建议
func (c *WorshipController) GetMessageReplySuggestions(ctx *gin.Context) {
	messageType := ctx.Query("message_type")
	content := ctx.Query("content")

	if messageType == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "留言类型不能为空",
		})
		return
	}

	suggestions, err := c.worshipService.GetMessageReplySuggestions(messageType, content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data: gin.H{
			"suggestions": suggestions,
		},
	})
}

// GetMessageCreationTips 获取留言创建提示
func (c *WorshipController) GetMessageCreationTips(ctx *gin.Context) {
	memorialID := ctx.Param("memorial_id")
	messageType := ctx.Query("message_type")

	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	if messageType == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "留言类型不能为空",
		})
		return
	}

	tips, err := c.worshipService.GetMessageCreationTips(memorialID, messageType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data: gin.H{
			"tips": tips,
		},
	})
}

// GetMemorialMessageAnalytics 获取纪念馆留言分析
func (c *WorshipController) GetMemorialMessageAnalytics(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	_ = userID // TODO: 应该使用 userID 进行权限验证

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	// 验证权限（只有纪念馆创建者可以查看分析数据）
	// 这里应该添加权限验证逻辑

	analytics, err := c.worshipService.GetMemorialMessageAnalytics(memorialID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    analytics,
	})
}

// GetDetailedWorshipStatistics 获取详细祭扫统计
func (c *WorshipController) GetDetailedWorshipStatistics(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	stats, err := c.worshipService.GetDetailedWorshipStatistics(userID.(string), memorialID)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    stats,
	})
}

// AnalyzeUserWorshipBehavior 分析用户祭扫行为
func (c *WorshipController) AnalyzeUserWorshipBehavior(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	behavior, err := c.worshipService.AnalyzeUserWorshipBehavior(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "分析完成",
		Data:    behavior,
	})
}

// GenerateWorshipReport 生成祭扫报告
func (c *WorshipController) GenerateWorshipReport(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	period := ctx.DefaultQuery("period", "month")
	validPeriods := []string{"week", "month", "quarter", "year"}
	isValidPeriod := false
	for _, validPeriod := range validPeriods {
		if period == validPeriod {
			isValidPeriod = true
			break
		}
	}
	if !isValidPeriod {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "无效的统计周期",
		})
		return
	}

	report, err := c.worshipService.GenerateWorshipReport(userID.(string), memorialID, period)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    3002,
				Message: err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, APIResponse{
				Code:    1005,
				Message: err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "报告生成成功",
		Data:    report,
	})
}
