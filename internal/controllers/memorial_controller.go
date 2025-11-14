package controllers

import (
	"net/http"
	"strconv"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type MemorialController struct {
	memorialService *services.MemorialService
}

func NewMemorialController(memorialService *services.MemorialService) *MemorialController {
	return &MemorialController{
		memorialService: memorialService,
	}
}

// CreateMemorial 创建纪念馆
func (c *MemorialController) CreateMemorial(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	var req services.CreateMemorialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	memorial, err := c.memorialService.CreateMemorial(userID.(string), &req)
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
		Data:    memorial,
	})
}

// GetMemorial 获取纪念馆详情
func (c *MemorialController) GetMemorial(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	memorial, err := c.memorialService.GetMemorial(userID.(string), memorialID)
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
		Data:    memorial,
	})
}

// UpdateMemorial 更新纪念馆信息
func (c *MemorialController) UpdateMemorial(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	var req services.UpdateMemorialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.memorialService.UpdateMemorial(userID.(string), memorialID, &req)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权修改此纪念馆" {
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
		Message: "更新成功",
	})
}

// DeleteMemorial 删除纪念馆
func (c *MemorialController) DeleteMemorial(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	err := c.memorialService.DeleteMemorial(userID.(string), memorialID)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权删除此纪念馆" {
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
		Message: "删除成功",
	})
}

// GetMemorialList 获取纪念馆列表
func (c *MemorialController) GetMemorialList(ctx *gin.Context) {
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

	memorials, total, err := c.memorialService.GetMemorialList(userID.(string), page, pageSize)
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
func (c *MemorialController) GetMemorialVisitors(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("id")
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

	visitors, total, err := c.memorialService.GetMemorialVisitors(userID.(string), memorialID, page, pageSize)
	if err != nil {
		if err.Error() == "只有创建者可以查看访客记录" {
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
			"list":      visitors,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetTombstoneStyles 获取墓碑样式列表
func (c *MemorialController) GetTombstoneStyles(ctx *gin.Context) {
	styles := c.memorialService.GetTombstoneStyles()

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    styles,
	})
}

// GetThemeStyles 获取主题风格列表
func (c *MemorialController) GetThemeStyles(ctx *gin.Context) {
	styles := c.memorialService.GetThemeStyles()

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "获取成功",
		Data:    styles,
	})
}

// UpdateTombstoneStyle 更新墓碑样式
func (c *MemorialController) UpdateTombstoneStyle(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	var req struct {
		StyleID string `json:"style_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.memorialService.UpdateTombstoneStyle(userID.(string), memorialID, req.StyleID)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权修改此纪念馆" {
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
		Message: "更新成功",
	})
}

// UpdateEpitaph 更新墓志铭
func (c *MemorialController) UpdateEpitaph(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	memorialID := ctx.Param("id")
	if memorialID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "纪念馆ID不能为空",
		})
		return
	}

	var req struct {
		Epitaph string `json:"epitaph"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.memorialService.UpdateEpitaph(userID.(string), memorialID, req.Epitaph)
	if err != nil {
		if err.Error() == "纪念馆不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    3001,
				Message: err.Error(),
			})
		} else if err.Error() == "无权修改此纪念馆" {
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
		Message: "更新成功",
	})
}

// GenerateCalligraphy 生成书法字体
func (c *MemorialController) GenerateCalligraphy(ctx *gin.Context) {
	var req struct {
		Text      string `json:"text" binding:"required"`
		FontStyle string `json:"font_style"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	result, err := c.memorialService.GenerateCalligraphy(req.Text, req.FontStyle)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "生成成功",
		Data:    result,
	})
}

// ProcessHandwriting 处理手写照片
func (c *MemorialController) ProcessHandwriting(ctx *gin.Context) {
	var req struct {
		ImageURL string `json:"image_url" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	result, err := c.memorialService.ProcessHandwritingImage(req.ImageURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    1005,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "处理成功",
		Data:    result,
	})
}

// GetRecentMemorials 获取最近访问的纪念馆
func (c *MemorialController) GetRecentMemorials(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	// 获取限制参数
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))
	if limit < 1 || limit > 20 {
		limit = 5
	}

	memorials, err := c.memorialService.GetRecentMemorials(userID.(string), limit)
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
		Data:    memorials,
	})
}
