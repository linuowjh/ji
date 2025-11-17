package controllers

import (
	"net/http"
	"strconv"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// WechatLogin 微信小程序登录
func (c *UserController) WechatLogin(ctx *gin.Context) {
	var req services.WechatLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.userService.WechatLogin(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "登录成功",
		Data:    resp,
	})
}

// GetUserInfo 获取用户信息
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	user, err := c.userService.GetUserInfo(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, APIResponse{
			Code:    2001,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    user,
	})
}

// UpdateUserInfo 更新用户信息
func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.userService.UpdateUserInfo(userID.(string), req.Nickname, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "更新成功",
	})
}

// GetUserMemorials 获取用户纪念馆列表
func (c *UserController) GetUserMemorials(ctx *gin.Context) {
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

	memorials, total, err := c.userService.GetUserMemorials(userID.(string), page, pageSize)
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
			"list":      memorials,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUserWorshipRecords 获取用户祭扫记录
func (c *UserController) GetUserWorshipRecords(ctx *gin.Context) {
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

	records, total, err := c.userService.GetUserWorshipRecords(userID.(string), page, pageSize)
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
			"list":      records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUserMemorialDetails 获取用户纪念馆详细信息
func (c *UserController) GetUserMemorialDetails(ctx *gin.Context) {
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

	memorials, total, err := c.userService.GetUserMemorialDetails(userID.(string), page, pageSize)
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
			"list":      memorials,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetMemorialVisitors 获取纪念馆访客记录
func (c *UserController) GetMemorialVisitors(ctx *gin.Context) {
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
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	visitors, total, err := c.userService.GetMemorialVisitors(userID.(string), memorialID, page, pageSize)
	if err != nil {
		if err.Error() == "纪念馆不存在或无权访问" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
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
			"list":      visitors,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUserFamilies 获取用户参与的家族圈
func (c *UserController) GetUserFamilies(ctx *gin.Context) {
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

	families, total, err := c.userService.GetUserFamilies(userID.(string), page, pageSize)
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
			"list":      families,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUserStatistics 获取用户统计信息
func (c *UserController) GetUserStatistics(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	stats, err := c.userService.GetUserStatistics(userID.(string))
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
		Data:    stats,
	})
}

// GetUserRecentActivities 获取用户最近活动
func (c *UserController) GetUserRecentActivities(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	// 获取限制参数
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	activities, err := c.userService.GetUserRecentActivities(userID.(string), limit)
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
		Data:    activities,
	})
}

// GetUserDashboard 获取用户仪表板数据
func (c *UserController) GetUserDashboard(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	// 获取用户统计信息
	stats, err := c.userService.GetUserStatistics(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取统计信息失败: " + err.Error(),
		})
		return
	}

	// 获取最近活动
	activities, err := c.userService.GetUserRecentActivities(userID.(string), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取最近活动失败: " + err.Error(),
		})
		return
	}

	// 获取用户纪念馆（前5个）
	memorials, _, err := c.userService.GetUserMemorials(userID.(string), 1, 5)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取纪念馆列表失败: " + err.Error(),
		})
		return
	}

	// 获取用户家族圈（前5个）
	families, _, err := c.userService.GetUserFamilies(userID.(string), 1, 5)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取家族圈列表失败: " + err.Error(),
		})
		return
	}

	dashboard := gin.H{
		"statistics":        stats,
		"recent_activities": activities,
		"memorials":         memorials,
		"families":          families,
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    dashboard,
	})
}

// GetUpcomingReminders 获取用户所有家族的即将到来的纪念日提醒
func (c *UserController) GetUpcomingReminders(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	reminders, err := c.userService.GetUpcomingReminders(userID.(string))
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
		Data:    reminders,
	})
}

// UpdatePhone 更新用户手机号
func (c *UserController) UpdatePhone(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	var req services.UpdatePhoneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	phone, err := c.userService.UpdatePhone(userID.(string), req.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "更新成功",
		Data: gin.H{
			"phone": phone,
		},
	})
}
