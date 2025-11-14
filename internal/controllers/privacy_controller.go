package controllers

import (
	"net/http"
	"strconv"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type PrivacyController struct {
	privacyService *services.PrivacyService
}

func NewPrivacyController(privacyService *services.PrivacyService) *PrivacyController {
	return &PrivacyController{
		privacyService: privacyService,
	}
}

// SetMemorialPrivacy 设置纪念馆隐私
func (c *PrivacyController) SetMemorialPrivacy(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	var req services.PrivacySettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.privacyService.SetMemorialPrivacy(userID.(string), &req)
	if err != nil {
		if err.Error() == "纪念馆不存在或无权操作" {
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
		Message: "隐私设置成功",
	})
}

// GetMemorialPrivacySettings 获取纪念馆隐私设置
func (c *PrivacyController) GetMemorialPrivacySettings(ctx *gin.Context) {
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

	settings, err := c.privacyService.GetMemorialPrivacySettings(userID.(string), memorialID)
	if err != nil {
		if err.Error() == "纪念馆不存在或无权操作" {
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
		Data:    settings,
	})
}

// CheckUserAccess 检查用户访问权限
func (c *PrivacyController) CheckUserAccess(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	permissionType := ctx.Query("permission_type")
	
	if memorialID == "" || permissionType == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID和权限类型不能为空",
		})
		return
	}

	hasAccess, err := c.privacyService.CheckUserAccess(userID.(string), memorialID, permissionType)
	if err != nil {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "检查完成",
		Data: gin.H{
			"has_access": hasAccess,
		},
	})
}

// RequestAccess 申请访问权限
func (c *PrivacyController) RequestAccess(ctx *gin.Context) {
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
		Message string `json:"message"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.privacyService.RequestAccess(userID.(string), memorialID, req.Message)
	if err != nil {
		if err.Error() == "已有待处理的访问申请" {
			ctx.JSON(http.StatusBadRequest, APIResponse{
				Code:    1001,
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
		Message: "申请已提交",
	})
}

// HandleAccessRequest 处理访问申请
func (c *PrivacyController) HandleAccessRequest(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	requestID := ctx.Param("request_id")
	if requestID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "申请ID不能为空",
		})
		return
	}

	var req struct {
		Approve     bool     `json:"approve"`
		Permissions []string `json:"permissions"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.privacyService.HandleAccessRequest(userID.(string), requestID, req.Approve, req.Permissions)
	if err != nil {
		if err.Error() == "访问申请不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "无权处理此申请" {
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

	message := "申请已拒绝"
	if req.Approve {
		message = "申请已批准"
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: message,
	})
}

// AddToBlacklist 添加用户到黑名单
func (c *PrivacyController) AddToBlacklist(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	targetUserID := ctx.Param("user_id")
	
	if memorialID == "" || targetUserID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID和用户ID不能为空",
		})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.privacyService.AddToBlacklist(userID.(string), memorialID, targetUserID, req.Reason)
	if err != nil {
		if err.Error() == "纪念馆不存在或无权操作" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "用户已在黑名单中" {
			ctx.JSON(http.StatusBadRequest, APIResponse{
				Code:    1001,
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
		Message: "已添加到黑名单",
	})
}

// RemoveFromBlacklist 从黑名单移除用户
func (c *PrivacyController) RemoveFromBlacklist(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("memorial_id")
	targetUserID := ctx.Param("user_id")
	
	if memorialID == "" || targetUserID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID和用户ID不能为空",
		})
		return
	}

	err := c.privacyService.RemoveFromBlacklist(userID.(string), memorialID, targetUserID)
	if err != nil {
		if err.Error() == "纪念馆不存在或无权操作" {
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
		Message: "已从黑名单移除",
	})
}

// GetAccessRequests 获取访问申请列表
func (c *PrivacyController) GetAccessRequests(ctx *gin.Context) {
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

	requests, total, err := c.privacyService.GetAccessRequests(userID.(string), memorialID, page, pageSize)
	if err != nil {
		if err.Error() == "纪念馆不存在或无权操作" {
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
			"list":      requests,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}