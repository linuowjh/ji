package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminService *services.AdminService
}

func NewAdminController(adminService *services.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

// GetUserList 获取用户列表
func (c *AdminController) GetUserList(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	// 获取查询参数
	req := &services.UserSearchRequest{
		Keyword:   ctx.Query("keyword"),
		Status:    -1, // 默认查询所有状态
		Role:      ctx.Query("role"),
		StartDate: ctx.Query("start_date"),
		EndDate:   ctx.Query("end_date"),
		Page:      1,
		PageSize:  20,
	}

	// 解析状态参数
	if statusStr := ctx.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			req.Status = status
		}
	}

	// 解析分页参数
	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}
	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			req.PageSize = pageSize
		}
	}

	users, total, err := c.adminService.GetUserList(userID.(string), req)
	if err != nil {
		if err.Error() == "权限不足" {
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
			"list":      users,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// GetUserDetail 获取用户详细信息
func (c *AdminController) GetUserDetail(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	targetUserID := ctx.Param("user_id")
	if targetUserID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "用户ID不能为空",
		})
		return
	}

	userDetail, err := c.adminService.GetUserDetail(userID.(string), targetUserID)
	if err != nil {
		if err.Error() == "权限不足" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "用户不存在" {
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
		Message: "获取成功",
		Data:    userDetail,
	})
}

// ManageUser 管理用户状态
func (c *AdminController) ManageUser(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	var req services.UserManagementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.adminService.ManageUser(userID.(string), &req)
	if err != nil {
		if err.Error() == "权限不足" || err.Error() == "无权操作超级管理员" || err.Error() == "不能操作自己的账户" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "用户不存在" || err.Error() == "无效的操作类型" {
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

	actionMessages := map[string]string{
		"activate":   "用户已激活",
		"deactivate": "用户已禁用",
		"approve":    "用户已审核通过",
		"reject":     "用户已审核拒绝",
	}

	message := actionMessages[req.Action]
	if message == "" {
		message = "操作成功"
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: message,
	})
}

// GetSystemStats 获取系统统计信息
func (c *AdminController) GetSystemStats(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	stats, err := c.adminService.GetSystemStats(userID.(string))
	if err != nil {
		if err.Error() == "权限不足" {
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
		Data:    stats,
	})
}

// GetPendingContent 获取待审核内容列表
func (c *AdminController) GetPendingContent(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	contentType := ctx.Query("content_type")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	contents, total, err := c.adminService.GetPendingContent(userID.(string), contentType, page, pageSize)
	if err != nil {
		if err.Error() == "权限不足" {
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
			"list":         contents,
			"total":        total,
			"page":         page,
			"page_size":    pageSize,
			"content_type": contentType,
		},
	})
}

// ModerateContent 审核内容
func (c *AdminController) ModerateContent(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	var req services.ContentModerationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.adminService.ModerateContent(userID.(string), &req)
	if err != nil {
		if err.Error() == "权限不足" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if strings.Contains(err.Error(), "不存在") || strings.Contains(err.Error(), "无效") || strings.Contains(err.Error(), "不支持") {
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

	actionMessages := map[string]string{
		"approve": "内容已审核通过",
		"reject":  "内容已审核拒绝",
	}

	message := actionMessages[req.Action]
	if message == "" {
		message = "审核操作成功"
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: message,
	})
}

// BatchModerateContent 批量审核内容
func (c *AdminController) BatchModerateContent(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	var req struct {
		ContentIDs  []string `json:"content_ids" binding:"required"`
		ContentType string   `json:"content_type" binding:"required,oneof=memorial message prayer"`
		Action      string   `json:"action" binding:"required,oneof=approve reject"`
		Reason      string   `json:"reason"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.adminService.BatchModerateContent(userID.(string), req.ContentIDs, req.ContentType, req.Action, req.Reason)
	if err != nil {
		if err.Error() == "权限不足" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if strings.Contains(err.Error(), "无效") || strings.Contains(err.Error(), "不能为空") || strings.Contains(err.Error(), "不支持") {
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

	actionMessages := map[string]string{
		"approve": "批量审核通过成功",
		"reject":  "批量审核拒绝成功",
	}

	message := actionMessages[req.Action]
	if message == "" {
		message = "批量审核操作成功"
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: message,
	})
}