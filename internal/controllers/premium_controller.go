package controllers

import (
	"net/http"
	"strconv"
	"time"
	"yun-nian-memorial/internal/models"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type PremiumController struct {
	premiumService *services.PremiumService
}

func NewPremiumController(premiumService *services.PremiumService) *PremiumController {
	return &PremiumController{
		premiumService: premiumService,
	}
}

// GetPremiumPackages 获取高级套餐列表
func (c *PremiumController) GetPremiumPackages(ctx *gin.Context) {
	packageType := ctx.Query("package_type")
	activeOnly := ctx.DefaultQuery("active_only", "true") == "true"
	
	packages, err := c.premiumService.GetPremiumPackages(packageType, activeOnly)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取套餐列表失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    packages,
	})
}

// GetPremiumPackage 获取套餐详情
func (c *PremiumController) GetPremiumPackage(ctx *gin.Context) {
	packageID := ctx.Param("id")
	
	pkg, err := c.premiumService.GetPremiumPackage(packageID)
	if err != nil {
		if err.Error() == "套餐不存在" {
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
		Data:    pkg,
	})
}

// Subscribe 订阅套餐
func (c *PremiumController) Subscribe(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	var req struct {
		PackageID  string `json:"package_id" binding:"required"`
		MemorialID string `json:"memorial_id"`
	}
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	subscription, err := c.premiumService.Subscribe(userID.(string), req.PackageID, req.MemorialID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "订阅失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "订阅成功",
		Data:    subscription,
	})
}

// GetUserSubscriptions 获取用户订阅列表
func (c *PremiumController) GetUserSubscriptions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	status := ctx.Query("status")
	
	subscriptions, err := c.premiumService.GetUserSubscriptions(userID.(string), status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取订阅列表失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    subscriptions,
	})
}

// GetSubscription 获取订阅详情
func (c *PremiumController) GetSubscription(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	subscriptionID := ctx.Param("id")
	
	subscription, err := c.premiumService.GetSubscription(subscriptionID)
	if err != nil {
		if err.Error() == "订阅不存在" {
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
	
	// 验证订阅属于当前用户
	if subscription.UserID != userID.(string) {
		ctx.JSON(http.StatusForbidden, APIResponse{
			Code:    1003,
			Message: "无权访问",
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    subscription,
	})
}

// CancelSubscription 取消订阅
func (c *PremiumController) CancelSubscription(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	subscriptionID := ctx.Param("id")
	
	if err := c.premiumService.CancelSubscription(subscriptionID, userID.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "取消订阅失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "取消成功",
	})
}

// RenewSubscription 续订
func (c *PremiumController) RenewSubscription(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	subscriptionID := ctx.Param("id")
	
	if err := c.premiumService.RenewSubscription(subscriptionID, userID.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "续订失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "续订成功",
	})
}

// UpgradeMemorial 升级纪念馆
func (c *PremiumController) UpgradeMemorial(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	_ = userID // TODO: 应该使用 userID 进行权限验证
	
	var req struct {
		MemorialID     string                 `json:"memorial_id" binding:"required"`
		SubscriptionID string                 `json:"subscription_id" binding:"required"`
		UpgradeType    string                 `json:"upgrade_type" binding:"required"`
		UpgradeData    map[string]interface{} `json:"upgrade_data"`
	}
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	if err := c.premiumService.UpgradeMemorial(req.MemorialID, req.SubscriptionID, req.UpgradeType, req.UpgradeData); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "升级失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "升级成功",
	})
}

// GetMemorialUpgrades 获取纪念馆升级记录
func (c *PremiumController) GetMemorialUpgrades(ctx *gin.Context) {
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
	
	upgrades, err := c.premiumService.GetMemorialUpgrades(memorialID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取升级记录失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    upgrades,
	})
}

// CreateCustomTemplate 创建定制模板
func (c *PremiumController) CreateCustomTemplate(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	var template models.CustomTemplate
	if err := ctx.ShouldBindJSON(&template); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	template.UserID = userID.(string)
	
	if err := c.premiumService.CreateCustomTemplate(&template); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "创建模板失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "创建成功",
		Data:    template,
	})
}

// GetUserCustomTemplates 获取用户定制模板列表
func (c *PremiumController) GetUserCustomTemplates(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	templateType := ctx.Query("template_type")
	
	templates, err := c.premiumService.GetUserCustomTemplates(userID.(string), templateType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取模板列表失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    templates,
	})
}

// UpdateCustomTemplate 更新定制模板
func (c *PremiumController) UpdateCustomTemplate(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
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
	
	if err := c.premiumService.UpdateCustomTemplate(templateID, userID.(string), updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "更新模板失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "更新成功",
	})
}

// DeleteCustomTemplate 删除定制模板
func (c *PremiumController) DeleteCustomTemplate(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	templateID := ctx.Param("id")
	
	if err := c.premiumService.DeleteCustomTemplate(templateID, userID.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "删除模板失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "删除成功",
	})
}

// GetUserStorage 获取用户存储使用情况
func (c *PremiumController) GetUserStorage(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	storage, err := c.premiumService.GetUserStorage(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取存储信息失败: " + err.Error(),
		})
		return
	}
	
	// 计算使用百分比
	usagePercent := float64(0)
	if storage.TotalSpace > 0 {
		usagePercent = float64(storage.UsedSpace) / float64(storage.TotalSpace) * 100
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data: gin.H{
			"storage":        storage,
			"usage_percent":  usagePercent,
			"available_space": storage.TotalSpace - storage.UsedSpace,
		},
	})
}

// GetServiceUsageStats 获取服务使用统计
func (c *PremiumController) GetServiceUsageStats(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	
	// 解析时间范围
	daysAgo, _ := strconv.Atoi(ctx.DefaultQuery("days_ago", "30"))
	startTime := ctx.Query("start_time")
	endTime := ctx.Query("end_time")
	
	var start, end time.Time
	if startTime != "" {
		start, _ = time.Parse(time.RFC3339, startTime)
	} else {
		start = time.Now().AddDate(0, 0, -daysAgo)
	}
	
	if endTime != "" {
		end, _ = time.Parse(time.RFC3339, endTime)
	} else {
		end = time.Now()
	}
	
	stats, err := c.premiumService.GetServiceUsageStats(userID.(string), start, end)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "获取使用统计失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    stats,
	})
}

// CreatePremiumPackage 创建高级套餐（管理员）
func (c *PremiumController) CreatePremiumPackage(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	_ = userID // TODO: 检查管理员权限
	
	var pkg models.PremiumPackage
	if err := ctx.ShouldBindJSON(&pkg); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	if err := c.premiumService.CreatePremiumPackage(&pkg); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "创建套餐失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "创建成功",
		Data:    pkg,
	})
}

// UpdatePremiumPackage 更新套餐信息（管理员）
func (c *PremiumController) UpdatePremiumPackage(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	_ = userID // TODO: 检查管理员权限
	
	packageID := ctx.Param("id")
	
	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}
	
	if err := c.premiumService.UpdatePremiumPackage(packageID, updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "更新套餐失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "更新成功",
	})
}

// InitDefaultPackages 初始化默认套餐（管理员）
func (c *PremiumController) InitDefaultPackages(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}
	_ = userID // TODO: 检查管理员权限
	
	if err := c.premiumService.InitDefaultPackages(); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: "初始化默认套餐失败: " + err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "初始化成功",
	})
}
