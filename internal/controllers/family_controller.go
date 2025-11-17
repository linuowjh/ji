package controllers

import (
	"net/http"
	"strconv"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type FamilyController struct {
	familyService *services.FamilyService
}

func NewFamilyController(familyService *services.FamilyService) *FamilyController {
	return &FamilyController{
		familyService: familyService,
	}
}

// CreateFamily 创建家族圈
func (c *FamilyController) CreateFamily(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	var req services.CreateFamilyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	family, err := c.familyService.CreateFamily(userID.(string), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "创建成功",
		Data:    family,
	})
}

// GetFamilies 获取家族圈列表
func (c *FamilyController) GetFamilies(ctx *gin.Context) {
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

	families, total, err := c.familyService.GetFamilies(userID.(string), page, pageSize)
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

// GetFamily 获取家族圈详情
func (c *FamilyController) GetFamily(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	family, err := c.familyService.GetFamily(userID.(string), familyID)
	if err != nil {
		if err.Error() == "家族圈不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "您不是此家族圈的成员" {
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

	// 获取统计数据
	stats, err := c.familyService.GetFamilyStatistics(familyID)
	if err != nil {
		// 统计数据获取失败不影响主流程，使用默认值
		stats = map[string]interface{}{
			"memberCount":   0,
			"memorialCount": 0,
			"activityCount": 0,
		}
	}

	// 合并family数据和统计数据
	response := gin.H{
		"id":            family.ID,
		"name":          family.Name,
		"creatorId":     family.CreatorID,
		"description":   family.Description,
		"inviteCode":    family.InviteCode,
		"createdAt":     family.CreatedAt,
		"updatedAt":     family.UpdatedAt,
		"creator":       family.Creator,
		"memberCount":   stats["memberCount"],
		"memorialCount": stats["memorialCount"],
		"activityCount": stats["activityCount"],
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    response,
	})
}

// UpdateFamily 更新家族圈
func (c *FamilyController) UpdateFamily(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	var req services.UpdateFamilyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.familyService.UpdateFamily(userID.(string), familyID, &req)
	if err != nil {
		if err.Error() == "家族圈不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "只有管理员可以修改家族圈" {
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
		Message: "更新成功",
	})
}

// DeleteFamily 删除家族圈
func (c *FamilyController) DeleteFamily(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	err := c.familyService.DeleteFamily(userID.(string), familyID)
	if err != nil {
		if err.Error() == "家族圈不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "只有创建者可以删除家族圈" {
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
		Message: "删除成功",
	})
}

// InviteMembers 邀请成员
func (c *FamilyController) InviteMembers(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	var req services.InviteFamilyMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.familyService.InviteMembers(userID.(string), familyID, &req)
	if err != nil {
		if err.Error() == "只有管理员可以邀请成员" {
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
		Message: "邀请发送成功",
	})
}

// JoinFamilyByCode 通过邀请码加入家族圈
func (c *FamilyController) JoinFamilyByCode(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	var req struct {
		InviteCode string `json:"inviteCode" binding:"required"` // 驼峰命名
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.familyService.JoinFamilyByCode(userID.(string), req.InviteCode)
	if err != nil {
		if err.Error() == "邀请码无效" {
			ctx.JSON(http.StatusBadRequest, APIResponse{
				Code:    1001,
				Message: err.Error(),
			})
		} else if err.Error() == "您已经是此家族圈的成员" {
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
		Message: "加入成功",
	})
}

// RespondToInvitation 响应邀请
func (c *FamilyController) RespondToInvitation(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	invitationID := ctx.Param("invitation_id")
	if invitationID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "邀请ID不能为空",
		})
		return
	}

	var req struct {
		Accept bool `json:"accept"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.familyService.RespondToInvitation(userID.(string), invitationID, req.Accept)
	if err != nil {
		if err.Error() == "邀请不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
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

	message := "邀请已拒绝"
	if req.Accept {
		message = "邀请已接受"
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: message,
	})
}

// RemoveMember 移除成员
func (c *FamilyController) RemoveMember(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	memberID := ctx.Param("member_id")
	if familyID == "" || memberID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和成员ID不能为空",
		})
		return
	}

	err := c.familyService.RemoveMember(userID.(string), familyID, memberID)
	if err != nil {
		if err.Error() == "只有管理员可以移除成员" || err.Error() == "不能移除家族圈创建者" {
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
		Message: "移除成功",
	})
}

// LeaveFamily 离开家族圈
func (c *FamilyController) LeaveFamily(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	err := c.familyService.LeaveFamily(userID.(string), familyID)
	if err != nil {
		if err.Error() == "创建者不能离开家族圈" {
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
		Message: "离开成功",
	})
}

// SetMemberRole 设置成员角色
func (c *FamilyController) SetMemberRole(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	memberID := ctx.Param("member_id")
	if familyID == "" || memberID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和成员ID不能为空",
		})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.familyService.SetMemberRole(userID.(string), familyID, memberID, req.Role)
	if err != nil {
		if err.Error() == "只有创建者可以设置成员角色" || err.Error() == "不能修改创建者的角色" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "无效的角色" {
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
		Message: "设置成功",
	})
}

// GetFamilyMembers 获取家族成员列表
func (c *FamilyController) GetFamilyMembers(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
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

	members, total, err := c.familyService.GetFamilyMembers(userID.(string), familyID, page, pageSize)
	if err != nil {
		if err.Error() == "您不是此家族圈的成员" {
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
			"list":      members,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetFamilyActivities 获取家族活动
func (c *FamilyController) GetFamilyActivities(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
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

	activities, total, err := c.familyService.GetFamilyActivities(userID.(string), familyID, page, pageSize)
	if err != nil {
		if err.Error() == "您不是此家族圈的成员" {
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
			"list":      activities,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AddMemorialToFamily 关联纪念馆到家族圈
func (c *FamilyController) AddMemorialToFamily(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	var req struct {
		MemorialID string `json:"memorial_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.familyService.AddMemorialToFamily(userID.(string), familyID, req.MemorialID)
	if err != nil {
		if err.Error() == "只有管理员可以关联纪念馆" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "纪念馆不存在或无权操作" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "纪念馆已经关联到此家族圈" {
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
		Message: "关联成功",
	})
}

// RemoveMemorialFromFamily 移除纪念馆关联
func (c *FamilyController) RemoveMemorialFromFamily(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	memorialID := ctx.Param("memorial_id")
	if familyID == "" || memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和纪念馆ID不能为空",
		})
		return
	}

	err := c.familyService.RemoveMemorialFromFamily(userID.(string), familyID, memorialID)
	if err != nil {
		if err.Error() == "只有管理员可以移除纪念馆关联" {
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
		Message: "移除成功",
	})
}

// SetMemorialReminder 设置纪念日提醒
func (c *FamilyController) SetMemorialReminder(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	var req services.SetReminderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.familyService.SetMemorialReminder(userID.(string), familyID, &req)
	if err != nil {
		if err.Error() == "只有管理员可以设置纪念日提醒" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "纪念馆未关联到此家族圈" || err.Error() == "无效的提醒类型" {
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
		Message: "设置成功",
	})
}

// GetFamilyReminders 获取家族纪念日提醒
func (c *FamilyController) GetFamilyReminders(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
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

	reminders, total, err := c.familyService.GetFamilyReminders(userID.(string), familyID, page, pageSize)
	if err != nil {
		if err.Error() == "您不是此家族圈的成员" {
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
			"list":      reminders,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUpcomingReminders 获取即将到来的提醒
func (c *FamilyController) GetUpcomingReminders(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	reminders, err := c.familyService.GetUpcomingReminders(userID.(string), familyID)
	if err != nil {
		if err.Error() == "您不是此家族圈的成员" {
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
		Data:    reminders,
	})
}

// DeleteReminder 删除纪念日提醒
func (c *FamilyController) DeleteReminder(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	reminderID := ctx.Param("reminder_id")
	if familyID == "" || reminderID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和提醒ID不能为空",
		})
		return
	}

	err := c.familyService.DeleteReminder(userID.(string), familyID, reminderID)
	if err != nil {
		if err.Error() == "只有管理员可以删除纪念日提醒" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "提醒不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "无权删除此提醒" {
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
		Message: "删除成功",
	})
}

// InitiateCollectiveWorship 发起集体祭扫
func (c *FamilyController) InitiateCollectiveWorship(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	var req services.CollectiveWorshipRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.familyService.InitiateCollectiveWorship(userID.(string), familyID, &req)
	if err != nil {
		if err.Error() == "只有管理员可以发起集体祭扫" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "纪念馆未关联到此家族圈" || err.Error() == "无效的祭扫类型" {
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
		Message: "集体祭扫发起成功",
	})
}

// JoinCollectiveWorship 参与集体祭扫
func (c *FamilyController) JoinCollectiveWorship(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	activityID := ctx.Param("activity_id")
	if familyID == "" || activityID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和活动ID不能为空",
		})
		return
	}

	err := c.familyService.JoinCollectiveWorship(userID.(string), familyID, activityID)
	if err != nil {
		if err.Error() == "您不是此家族圈的成员" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "集体祭扫活动不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "活动已结束" || err.Error() == "您已经参与了此活动" {
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
		Message: "参与成功",
	})
}

// CreateGenealogy 创建家族谱系成员
func (c *FamilyController) CreateGenealogy(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	var req services.CreateGenealogyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	genealogy, err := c.familyService.CreateGenealogy(userID.(string), familyID, &req)
	if err != nil {
		if err.Error() == "只有管理员可以创建家族谱系" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "无效的性别" || err.Error() == "指定的父辈不存在" || err.Error() == "纪念馆未关联到此家族圈" {
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
		Message: "创建成功",
		Data:    genealogy,
	})
}

// GetFamilyGenealogy 获取家族谱系
func (c *FamilyController) GetFamilyGenealogy(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	genealogies, err := c.familyService.GetFamilyGenealogy(userID.(string), familyID)
	if err != nil {
		if err.Error() == "您不是此家族圈的成员" {
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
		Data:    genealogies,
	})
}

// UpdateGenealogy 更新家族谱系成员
func (c *FamilyController) UpdateGenealogy(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	genealogyID := ctx.Param("genealogy_id")
	if familyID == "" || genealogyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和谱系ID不能为空",
		})
		return
	}

	var req services.UpdateGenealogyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.familyService.UpdateGenealogy(userID.(string), familyID, genealogyID, &req)
	if err != nil {
		if err.Error() == "只有管理员可以更新家族谱系" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "谱系成员不存在" {
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
		Message: "更新成功",
	})
}

// DeleteGenealogy 删除家族谱系成员
func (c *FamilyController) DeleteGenealogy(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	genealogyID := ctx.Param("genealogy_id")
	if familyID == "" || genealogyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和谱系ID不能为空",
		})
		return
	}

	err := c.familyService.DeleteGenealogy(userID.(string), familyID, genealogyID)
	if err != nil {
		if err.Error() == "只有管理员可以删除家族谱系" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "谱系成员不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "该成员有子代记录，无法删除" {
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
		Message: "删除成功",
	})
}

// CreateFamilyStory 创建家族故事
func (c *FamilyController) CreateFamilyStory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	var req services.CreateFamilyStoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	story, err := c.familyService.CreateFamilyStory(userID.(string), familyID, &req)
	if err != nil {
		if err.Error() == "您不是此家族圈的成员" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "无效的故事分类" {
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
		Message: "创建成功",
		Data:    story,
	})
}

// GetFamilyStories 获取家族故事列表
func (c *FamilyController) GetFamilyStories(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	// 获取查询参数
	category := ctx.Query("category")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	stories, total, err := c.familyService.GetFamilyStories(userID.(string), familyID, category, page, pageSize)
	if err != nil {
		if err.Error() == "您不是此家族圈的成员" {
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
			"list":      stories,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetFamilyStory 获取家族故事详情
func (c *FamilyController) GetFamilyStory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	storyID := ctx.Param("story_id")
	if familyID == "" || storyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和故事ID不能为空",
		})
		return
	}

	story, err := c.familyService.GetFamilyStory(userID.(string), familyID, storyID)
	if err != nil {
		if err.Error() == "您不是此家族圈的成员" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "故事不存在" {
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
		Data:    story,
	})
}

// UpdateFamilyStory 更新家族故事
func (c *FamilyController) UpdateFamilyStory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	storyID := ctx.Param("story_id")
	if familyID == "" || storyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和故事ID不能为空",
		})
		return
	}

	var req services.UpdateFamilyStoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.familyService.UpdateFamilyStory(userID.(string), familyID, storyID, &req)
	if err != nil {
		if err.Error() == "只有作者或管理员可以更新故事" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "故事不存在" {
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
		Message: "更新成功",
	})
}

// DeleteFamilyStory 删除家族故事
func (c *FamilyController) DeleteFamilyStory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	storyID := ctx.Param("story_id")
	if familyID == "" || storyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和故事ID不能为空",
		})
		return
	}

	err := c.familyService.DeleteFamilyStory(userID.(string), familyID, storyID)
	if err != nil {
		if err.Error() == "只有作者或管理员可以删除故事" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "故事不存在" {
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
		Message: "删除成功",
	})
}

// CreateFamilyTradition 创建家族传统
func (c *FamilyController) CreateFamilyTradition(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	var req services.CreateFamilyTraditionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	tradition, err := c.familyService.CreateFamilyTradition(userID.(string), familyID, &req)
	if err != nil {
		if err.Error() == "只有管理员可以创建家族传统" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "无效的传统分类" {
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
		Message: "创建成功",
		Data:    tradition,
	})
}

// GetFamilyTraditions 获取家族传统列表
func (c *FamilyController) GetFamilyTraditions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	if familyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID不能为空",
		})
		return
	}

	// 获取查询参数
	category := ctx.Query("category")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	traditions, total, err := c.familyService.GetFamilyTraditions(userID.(string), familyID, category, page, pageSize)
	if err != nil {
		if err.Error() == "您不是此家族圈的成员" {
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
			"list":      traditions,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateFamilyTradition 更新家族传统
func (c *FamilyController) UpdateFamilyTradition(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	traditionID := ctx.Param("tradition_id")
	if familyID == "" || traditionID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和传统ID不能为空",
		})
		return
	}

	var req services.UpdateFamilyTraditionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.familyService.UpdateFamilyTradition(userID.(string), familyID, traditionID, &req)
	if err != nil {
		if err.Error() == "只有管理员可以更新家族传统" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "传统不存在" {
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
		Message: "更新成功",
	})
}

// DeleteFamilyTradition 删除家族传统
func (c *FamilyController) DeleteFamilyTradition(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	familyID := ctx.Param("family_id")
	traditionID := ctx.Param("tradition_id")
	if familyID == "" || traditionID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "家族圈ID和传统ID不能为空",
		})
		return
	}

	err := c.familyService.DeleteFamilyTradition(userID.(string), familyID, traditionID)
	if err != nil {
		if err.Error() == "只有管理员可以删除家族传统" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
				Message: err.Error(),
			})
		} else if err.Error() == "传统不存在" {
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
		Message: "删除成功",
	})
}
