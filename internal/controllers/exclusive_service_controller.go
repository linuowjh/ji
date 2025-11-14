package controllers

import (
	"net/http"
	"strconv"
	"yun-nian-memorial/internal/models"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type ExclusiveServiceController struct {
	exclusiveService *services.ExclusiveServiceService
}

func NewExclusiveServiceController(exclusiveService *services.ExclusiveServiceService) *ExclusiveServiceController {
	return &ExclusiveServiceController{
		exclusiveService: exclusiveService,
	}
}

// GetExclusiveServices 获取专属服务列表
func (c *ExclusiveServiceController) GetExclusiveServices(ctx *gin.Context) {
	serviceType := ctx.Query("service_type")
	activeOnly := ctx.DefaultQuery("active_only", "true") == "true"
	
	services, err := c.exclusiveService.GetExclusiveServices(serviceType, activeOnly)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取服务列表失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    services,
	})
}

// GetExclusiveService 获取专属服务详情
func (c *ExclusiveServiceController) GetExclusiveService(ctx *gin.Context) {
	serviceID := ctx.Param("id")
	
	service, err := c.exclusiveService.GetExclusiveService(serviceID)
	if err != nil {
		if err.Error() == "服务不存在" {
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
		Data:    service,
	})
}

// CreateBooking 创建服务预订
func (c *ExclusiveServiceController) CreateBooking(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	var booking models.ServiceBooking
	if err := ctx.ShouldBindJSON(&booking); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	booking.UserID = userID.(string)
	
	if err := c.exclusiveService.CreateBooking(&booking); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "创建预订失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "预订成功",
		Data:    booking,
	})
}

// GetUserBookings 获取用户预订列表
func (c *ExclusiveServiceController) GetUserBookings(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	status := ctx.Query("status")
	
	bookings, err := c.exclusiveService.GetUserBookings(userID.(string), status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取预订列表失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    bookings,
	})
}

// GetBooking 获取预订详情
func (c *ExclusiveServiceController) GetBooking(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	bookingID := ctx.Param("id")
	
	booking, err := c.exclusiveService.GetBooking(bookingID)
	if err != nil {
		if err.Error() == "预订不存在" {
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
	
	// 验证预订属于当前用户
	if booking.UserID != userID.(string) {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "无权访问",
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    booking,
	})
}

// CancelBooking 取消预订
func (c *ExclusiveServiceController) CancelBooking(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	bookingID := ctx.Param("id")
	
	var req struct {
		Reason string `json:"reason"`
	}
	ctx.ShouldBindJSON(&req)
	
	if err := c.exclusiveService.CancelBooking(bookingID, userID.(string), req.Reason); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "取消预订失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "取消成功",
	})
}

// CreateDataExportRequest 创建数据导出请求
func (c *ExclusiveServiceController) CreateDataExportRequest(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	var req models.DataExportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	req.UserID = userID.(string)
	
	if err := c.exclusiveService.CreateDataExportRequest(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "创建导出请求失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "导出请求已创建，正在后台处理",
		Data:    req,
	})
}

// GetUserExportRequests 获取用户导出请求列表
func (c *ExclusiveServiceController) GetUserExportRequests(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	requests, err := c.exclusiveService.GetUserExportRequests(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取导出请求列表失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    requests,
	})
}

// GetExportRequest 获取导出请求详情
func (c *ExclusiveServiceController) GetExportRequest(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	requestID := ctx.Param("id")
	
	request, err := c.exclusiveService.GetExportRequest(requestID)
	if err != nil {
		if err.Error() == "导出请求不存在" {
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
	
	// 验证请求属于当前用户
	if request.UserID != userID.(string) {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "无权访问",
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    request,
	})
}

// DownloadExport 下载导出文件
func (c *ExclusiveServiceController) DownloadExport(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	requestID := ctx.Param("id")
	
	filePath, err := c.exclusiveService.DownloadExport(requestID, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "下载失败: " + err.Error(),
		})
		return
	}
	
	ctx.File(filePath)
}

// CreatePhotoRestoreRequest 创建老照片修复请求
func (c *ExclusiveServiceController) CreatePhotoRestoreRequest(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	var req models.PhotoRestoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	req.UserID = userID.(string)
	
	if err := c.exclusiveService.CreatePhotoRestoreRequest(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "创建修复请求失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "修复请求已创建，正在处理中",
		Data:    req,
	})
}

// GetUserPhotoRestoreRequests 获取用户照片修复请求列表
func (c *ExclusiveServiceController) GetUserPhotoRestoreRequests(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	requests, err := c.exclusiveService.GetUserPhotoRestoreRequests(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取修复请求列表失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    requests,
	})
}

// CreateServiceReview 创建服务评价
func (c *ExclusiveServiceController) CreateServiceReview(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	var review models.ServiceReview
	if err := ctx.ShouldBindJSON(&review); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	review.UserID = userID.(string)
	
	if err := c.exclusiveService.CreateServiceReview(&review); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "创建评价失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "评价成功",
		Data:    review,
	})
}

// GetServiceReviews 获取服务评价列表
func (c *ExclusiveServiceController) GetServiceReviews(ctx *gin.Context) {
	serviceID := ctx.Param("service_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	reviews, total, err := c.exclusiveService.GetServiceReviews(serviceID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取评价列表失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data: gin.H{
			"list":      reviews,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// InitDefaultServices 初始化默认专属服务（管理员）
func (c *ExclusiveServiceController) InitDefaultServices(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	_ = userID // TODO: 检查管理员权限
	
	if err := c.exclusiveService.InitDefaultServices(); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "初始化默认服务失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "初始化成功",
	})
}
