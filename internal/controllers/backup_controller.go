package controllers

import (
	"net/http"
	"strconv"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type BackupController struct {
	backupService *services.BackupService
	adminService  *services.AdminService
}

func NewBackupController(backupService *services.BackupService, adminService *services.AdminService) *BackupController {
	return &BackupController{
		backupService: backupService,
		adminService:  adminService,
	}
}

// CreateBackup 创建数据备份
func (c *BackupController) CreateBackup(ctx *gin.Context) {
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
	
	var req struct {
		BackupType string `json:"backup_type" binding:"required,oneof=full incremental user"`
	}
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	backup, err := c.backupService.CreateBackup(req.BackupType, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "创建备份失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "备份任务已创建，正在后台处理",
		Data:    backup,
	})
}

// GetBackupList 获取备份列表
func (c *BackupController) GetBackupList(ctx *gin.Context) {
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
	backupType := ctx.Query("backup_type")
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	backups, total, err := c.backupService.GetBackupList(page, pageSize, backupType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取备份列表失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data: gin.H{
			"list":      backups,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetBackup 获取备份详情
func (c *BackupController) GetBackup(ctx *gin.Context) {
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
	
	backupID := ctx.Param("id")
	
	backup, err := c.backupService.GetBackup(backupID)
	if err != nil {
		if err.Error() == "备份不存在" {
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
		Data:    backup,
	})
}

// DeleteBackup 删除备份
func (c *BackupController) DeleteBackup(ctx *gin.Context) {
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
	
	backupID := ctx.Param("id")
	
	if err := c.backupService.DeleteBackup(backupID); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "删除备份失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "删除成功",
	})
}

// DownloadBackup 下载备份文件
func (c *BackupController) DownloadBackup(ctx *gin.Context) {
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
	
	backupID := ctx.Param("id")
	
	filePath, err := c.backupService.DownloadBackup(backupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "下载备份失败: " + err.Error(),
		})
		return
	}
	
	ctx.File(filePath)
}

// RestoreBackup 恢复备份
func (c *BackupController) RestoreBackup(ctx *gin.Context) {
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
	
	backupID := ctx.Param("id")
	
	if err := c.backupService.RestoreBackup(backupID); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "恢复备份失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "恢复成功",
	})
}

// CleanOldBackups 清理旧备份
func (c *BackupController) CleanOldBackups(ctx *gin.Context) {
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
	
	keepCount, _ := strconv.Atoi(ctx.DefaultQuery("keep_count", "10"))
	
	if err := c.backupService.CleanOldBackups(keepCount); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "清理旧备份失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "清理成功",
	})
}

// GetBackupStats 获取备份统计信息
func (c *BackupController) GetBackupStats(ctx *gin.Context) {
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
	
	stats, err := c.backupService.GetBackupStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取备份统计失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    stats,
	})
}
