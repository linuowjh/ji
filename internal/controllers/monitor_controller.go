package controllers

import (
	"net/http"
	"strconv"
	"time"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type MonitorController struct {
	monitorService *services.MonitorService
	loggingService *services.LoggingService
	adminService   *services.AdminService
}

func NewMonitorController(monitorService *services.MonitorService, loggingService *services.LoggingService, adminService *services.AdminService) *MonitorController {
	return &MonitorController{
		monitorService: monitorService,
		loggingService: loggingService,
		adminService:   adminService,
	}
}

// GetSystemHealth 获取系统健康状态
func (c *MonitorController) GetSystemHealth(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	health, err := c.monitorService.GetSystemHealth()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取系统健康状态失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    health,
	})
}

// GetDashboardStats 获取仪表板统计数据
func (c *MonitorController) GetDashboardStats(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	stats, err := c.monitorService.GetDashboardStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取仪表板统计失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    stats,
	})
}

// GetMetrics 获取监控指标
func (c *MonitorController) GetMetrics(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	metricType := ctx.Query("metric_type")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	
	// 解析时间范围
	var startTime, endTime time.Time
	if startStr := ctx.Query("start_time"); startStr != "" {
		startTime, _ = time.Parse(time.RFC3339, startStr)
	}
	if endStr := ctx.Query("end_time"); endStr != "" {
		endTime, _ = time.Parse(time.RFC3339, endStr)
	}
	
	metrics, err := c.monitorService.GetMetrics(metricType, startTime, endTime, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取监控指标失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    metrics,
	})
}

// GetMetricStats 获取指标统计
func (c *MonitorController) GetMetricStats(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	metricType := ctx.Query("metric_type")
	durationHours, _ := strconv.Atoi(ctx.DefaultQuery("duration_hours", "24"))
	duration := time.Duration(durationHours) * time.Hour
	
	stats, err := c.monitorService.GetMetricStats(metricType, duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取指标统计失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    stats,
	})
}

// GetMetricTrend 获取指标趋势
func (c *MonitorController) GetMetricTrend(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	metricType := ctx.Query("metric_type")
	durationHours, _ := strconv.Atoi(ctx.DefaultQuery("duration_hours", "24"))
	intervalMinutes, _ := strconv.Atoi(ctx.DefaultQuery("interval_minutes", "60"))
	
	duration := time.Duration(durationHours) * time.Hour
	interval := time.Duration(intervalMinutes) * time.Minute
	
	trend, err := c.monitorService.GetMetricTrend(metricType, duration, interval)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取指标趋势失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    trend,
	})
}

// GetLogs 获取日志列表
func (c *MonitorController) GetLogs(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	filter := &services.LogFilter{
		LogLevel:  ctx.Query("log_level"),
		LogType:   ctx.Query("log_type"),
		UserID:    ctx.Query("user_id"),
		Action:    ctx.Query("action"),
		IPAddress: ctx.Query("ip_address"),
		StartTime: ctx.Query("start_time"),
		EndTime:   ctx.Query("end_time"),
		Page:      1,
		PageSize:  20,
	}
	
	if page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "20")); err == nil && pageSize > 0 && pageSize <= 100 {
		filter.PageSize = pageSize
	}
	
	logs, total, err := c.loggingService.GetLogs(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取日志失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data: gin.H{
			"list":      logs,
			"total":     total,
			"page":      filter.Page,
			"page_size": filter.PageSize,
		},
	})
}

// GetLogStats 获取日志统计
func (c *MonitorController) GetLogStats(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	startTime := ctx.Query("start_time")
	endTime := ctx.Query("end_time")
	
	stats, err := c.loggingService.GetLogStats(startTime, endTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取日志统计失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    stats,
	})
}

// SearchLogs 搜索日志
func (c *MonitorController) SearchLogs(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	keyword := ctx.Query("keyword")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	logs, total, err := c.loggingService.SearchLogs(keyword, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "搜索日志失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "搜索成功",
		Data: gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetSecurityAlerts 获取安全警报
func (c *MonitorController) GetSecurityAlerts(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	alerts, total, err := c.loggingService.GetSecurityAlerts(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取安全警报失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data: gin.H{
			"list":      alerts,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CleanOldLogs 清理旧日志
func (c *MonitorController) CleanOldLogs(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	daysToKeep, _ := strconv.Atoi(ctx.DefaultQuery("days_to_keep", "30"))
	
	if err := c.loggingService.CleanOldLogs(daysToKeep); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "清理旧日志失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "清理成功",
	})
}

// CleanOldMetrics 清理旧监控数据
func (c *MonitorController) CleanOldMetrics(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	daysToKeep, _ := strconv.Atoi(ctx.DefaultQuery("days_to_keep", "7"))
	
	if err := c.monitorService.CleanOldMetrics(daysToKeep); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "清理旧监控数据失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "清理成功",
	})
}

// ExportLogs 导出日志
func (c *MonitorController) ExportLogs(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 检查管理员权限
	isAdmin, _, err := c.adminService.CheckAdminPermission(userID.(string))
	if err != nil || !isAdmin {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "权限不足",
		})
		return
	}
	
	filter := &services.LogFilter{
		LogLevel:  ctx.Query("log_level"),
		LogType:   ctx.Query("log_type"),
		UserID:    ctx.Query("user_id"),
		Action:    ctx.Query("action"),
		IPAddress: ctx.Query("ip_address"),
		StartTime: ctx.Query("start_time"),
		EndTime:   ctx.Query("end_time"),
		Page:      1,
		PageSize:  10000, // 导出时获取更多数据
	}
	
	jsonData, err := c.loggingService.ExportLogs(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "导出日志失败: " + err.Error(),
		})
		return
	}
	
	ctx.Header("Content-Type", "application/json")
	ctx.Header("Content-Disposition", "attachment; filename=logs_export.json")
	ctx.String(http.StatusOK, jsonData)
}
