package controllers

import (
	"net/http"
	"strconv"
	"yun-nian-memorial/internal/models"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type SystemConfigController struct {
	configService *services.SystemConfigService
	adminService  *services.AdminService
}

func NewSystemConfigController(configService *services.SystemConfigService, adminService *services.AdminService) *SystemConfigController {
	return &SystemConfigController{
		configService: configService,
		adminService:  adminService,
	}
}

// GetFestivalConfigs 获取祭扫节日配置列表
func (c *SystemConfigController) GetFestivalConfigs(ctx *gin.Context) {
	activeOnly := ctx.DefaultQuery("active_only", "true") == "true"
	
	festivals, err := c.configService.GetFestivalConfigs(activeOnly)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取节日配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    festivals,
	})
}

// GetFestivalConfig 获取单个节日配置
func (c *SystemConfigController) GetFestivalConfig(ctx *gin.Context) {
	festivalID := ctx.Param("id")
	
	festival, err := c.configService.GetFestivalConfig(festivalID)
	if err != nil {
		if err.Error() == "节日配置不存在" {
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
		Data:    festival,
	})
}

// CreateFestivalConfig 创建祭扫节日配置
func (c *SystemConfigController) CreateFestivalConfig(ctx *gin.Context) {
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
	
	var festival models.FestivalConfig
	if err := ctx.ShouldBindJSON(&festival); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	if err := c.configService.CreateFestivalConfig(&festival); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "创建节日配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "创建成功",
		Data:    festival,
	})
}

// UpdateFestivalConfig 更新祭扫节日配置
func (c *SystemConfigController) UpdateFestivalConfig(ctx *gin.Context) {
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
	
	festivalID := ctx.Param("id")
	
	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	if err := c.configService.UpdateFestivalConfig(festivalID, updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "更新节日配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "更新成功",
	})
}

// DeleteFestivalConfig 删除祭扫节日配置
func (c *SystemConfigController) DeleteFestivalConfig(ctx *gin.Context) {
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
	
	festivalID := ctx.Param("id")
	
	if err := c.configService.DeleteFestivalConfig(festivalID); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "删除节日配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "删除成功",
	})
}

// GetTemplateConfigs 获取模板配置列表
func (c *SystemConfigController) GetTemplateConfigs(ctx *gin.Context) {
	templateType := ctx.Query("template_type")
	activeOnly := ctx.DefaultQuery("active_only", "true") == "true"
	
	templates, err := c.configService.GetTemplateConfigs(templateType, activeOnly)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取模板配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    templates,
	})
}

// CreateTemplateConfig 创建模板配置
func (c *SystemConfigController) CreateTemplateConfig(ctx *gin.Context) {
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
	
	var template models.TemplateConfig
	if err := ctx.ShouldBindJSON(&template); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	if err := c.configService.CreateTemplateConfig(&template); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "创建模板配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "创建成功",
		Data:    template,
	})
}

// UpdateTemplateConfig 更新模板配置
func (c *SystemConfigController) UpdateTemplateConfig(ctx *gin.Context) {
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
	
	templateID := ctx.Param("id")
	
	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	if err := c.configService.UpdateTemplateConfig(templateID, updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "更新模板配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "更新成功",
	})
}

// DeleteTemplateConfig 删除模板配置
func (c *SystemConfigController) DeleteTemplateConfig(ctx *gin.Context) {
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
	
	templateID := ctx.Param("id")
	
	if err := c.configService.DeleteTemplateConfig(templateID); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "删除模板配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "删除成功",
	})
}

// GetSystemConfigs 获取系统配置列表
func (c *SystemConfigController) GetSystemConfigs(ctx *gin.Context) {
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
	
	configType := ctx.Query("config_type")
	
	configs, err := c.configService.GetSystemConfigsByType(configType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取系统配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    configs,
	})
}

// SetSystemConfig 设置系统配置
func (c *SystemConfigController) SetSystemConfig(ctx *gin.Context) {
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
		ConfigKey   string `json:"config_key" binding:"required"`
		ConfigValue string `json:"config_value" binding:"required"`
		ConfigType  string `json:"config_type" binding:"required"`
		Description string `json:"description"`
	}
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	if err := c.configService.SetSystemConfig(req.ConfigKey, req.ConfigValue, req.ConfigType, req.Description); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "设置系统配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "设置成功",
	})
}

// InitDefaultConfigs 初始化默认配置
func (c *SystemConfigController) InitDefaultConfigs(ctx *gin.Context) {
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
	
	if err := c.configService.InitDefaultConfigs(); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "初始化默认配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "初始化成功",
	})
}

// ExportConfig 导出配置
func (c *SystemConfigController) ExportConfig(ctx *gin.Context) {
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
	
	configType := ctx.DefaultQuery("config_type", "all")
	
	jsonData, err := c.configService.ExportConfig(configType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "导出配置失败: " + err.Error(),
		})
		return
	}
	
	ctx.Header("Content-Type", "application/json")
	ctx.Header("Content-Disposition", "attachment; filename=config_export.json")
	ctx.String(http.StatusOK, jsonData)
}

// GetUpcomingFestivals 获取即将到来的节日
func (c *SystemConfigController) GetUpcomingFestivals(ctx *gin.Context) {
	daysAhead, _ := strconv.Atoi(ctx.DefaultQuery("days_ahead", "7"))
	
	festivals, err := c.configService.GetUpcomingFestivals(daysAhead)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取即将到来的节日失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    festivals,
	})
}
