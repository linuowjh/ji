package controllers

import (
	"net/http"
	"strconv"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type MemorialServiceController struct {
	memorialServiceService *services.MemorialServiceService
}

func NewMemorialServiceController(memorialServiceService *services.MemorialServiceService) *MemorialServiceController {
	return &MemorialServiceController{
		memorialServiceService: memorialServiceService,
	}
}

// CreateMemorialService 创建追思会
func (c *MemorialServiceController) CreateMemorialService(ctx *gin.Context) {
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

	var req services.CreateMemorialServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	service, err := c.memorialServiceService.CreateMemorialService(userID.(string), memorialID, &req)
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
		Message: "创建成功",
		Data:    service,
	})
}

// GetMemorialServices 获取追思会列表
func (c *MemorialServiceController) GetMemorialServices(ctx *gin.Context) {
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

	services, total, err := c.memorialServiceService.GetMemorialServices(userID.(string), memorialID, page, pageSize)
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
			"list":      services,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetMemorialService 获取追思会详情
func (c *MemorialServiceController) GetMemorialService(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	serviceID := ctx.Param("service_id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "追思会ID不能为空",
		})
		return
	}

	service, err := c.memorialServiceService.GetMemorialService(userID.(string), serviceID)
	if err != nil {
		if err.Error() == "追思会不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "无权访问此追思会" {
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
		Data:    service,
	})
}

// UpdateMemorialService 更新追思会
func (c *MemorialServiceController) UpdateMemorialService(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	serviceID := ctx.Param("service_id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "追思会ID不能为空",
		})
		return
	}

	var req services.UpdateMemorialServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.memorialServiceService.UpdateMemorialService(userID.(string), serviceID, &req)
	if err != nil {
		if err.Error() == "追思会不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "只有主持人可以修改追思会" || err.Error() == "只有未开始的追思会可以修改" {
			ctx.JSON(http.StatusForbidden, APIResponse{
				Code:    1003,
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
		Message: "更新成功",
	})
}

// DeleteMemorialService 删除追思会
func (c *MemorialServiceController) DeleteMemorialService(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	serviceID := ctx.Param("service_id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "追思会ID不能为空",
		})
		return
	}

	err := c.memorialServiceService.DeleteMemorialService(userID.(string), serviceID)
	if err != nil {
		if err.Error() == "追思会不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "只有主持人可以删除追思会" || err.Error() == "进行中的追思会不能删除" {
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

// InviteParticipants 邀请参与者
func (c *MemorialServiceController) InviteParticipants(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	serviceID := ctx.Param("service_id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "追思会ID不能为空",
		})
		return
	}

	var req services.InviteParticipantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.memorialServiceService.InviteParticipants(userID.(string), serviceID, &req)
	if err != nil {
		if err.Error() == "追思会不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "无权邀请参与者" {
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

// RespondToInvitation 响应邀请
func (c *MemorialServiceController) RespondToInvitation(ctx *gin.Context) {
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

	err := c.memorialServiceService.RespondToInvitation(userID.(string), invitationID, req.Accept)
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

// JoinService 加入追思会
func (c *MemorialServiceController) JoinService(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	serviceID := ctx.Param("service_id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "追思会ID不能为空",
		})
		return
	}

	err := c.memorialServiceService.JoinService(userID.(string), serviceID)
	if err != nil {
		if err.Error() == "追思会不存在" {
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

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "加入成功",
	})
}

// LeaveService 离开追思会
func (c *MemorialServiceController) LeaveService(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	serviceID := ctx.Param("service_id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "追思会ID不能为空",
		})
		return
	}

	err := c.memorialServiceService.LeaveService(userID.(string), serviceID)
	if err != nil {
		if err.Error() == "您未参加此追思会" {
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

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "离开成功",
	})
}

// StartService 开始追思会
func (c *MemorialServiceController) StartService(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	serviceID := ctx.Param("service_id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "追思会ID不能为空",
		})
		return
	}

	err := c.memorialServiceService.StartService(userID.(string), serviceID)
	if err != nil {
		if err.Error() == "追思会不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "只有主持人可以开始追思会" || err.Error() == "追思会已开始或已结束" {
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
		Message: "追思会已开始",
	})
}

// EndService 结束追思会
func (c *MemorialServiceController) EndService(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	serviceID := ctx.Param("service_id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "追思会ID不能为空",
		})
		return
	}

	err := c.memorialServiceService.EndService(userID.(string), serviceID)
	if err != nil {
		if err.Error() == "追思会不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "只有主持人可以结束追思会" || err.Error() == "追思会未开始或已结束" {
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
		Message: "追思会已结束，正在生成录制视频",
	})
}

// SendChatMessage 发送聊天消息
func (c *MemorialServiceController) SendChatMessage(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	serviceID := ctx.Param("service_id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "追思会ID不能为空",
		})
		return
	}

	var req services.SendChatMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	message, err := c.memorialServiceService.SendChatMessage(userID.(string), serviceID, &req)
	if err != nil {
		if err.Error() == "您不是此追思会的参与者" {
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
		Message: "发送成功",
		Data:    message,
	})
}

// GetChatMessages 获取聊天消息
func (c *MemorialServiceController) GetChatMessages(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	serviceID := ctx.Param("service_id")
	if serviceID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "追思会ID不能为空",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	messages, total, err := c.memorialServiceService.GetChatMessages(userID.(string), serviceID, page, pageSize)
	if err != nil {
		if err.Error() == "您不是此追思会的参与者" {
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
			"list":      messages,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}