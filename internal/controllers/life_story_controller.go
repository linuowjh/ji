package controllers

import (
	"net/http"
	"strconv"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type LifeStoryController struct {
	lifeStoryService *services.LifeStoryService
}

func NewLifeStoryController(lifeStoryService *services.LifeStoryService) *LifeStoryController {
	return &LifeStoryController{
		lifeStoryService: lifeStoryService,
	}
}

// CreateLifeStory 创建生平故事
func (c *LifeStoryController) CreateLifeStory(ctx *gin.Context) {
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

	var req services.CreateLifeStoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	story, err := c.lifeStoryService.CreateLifeStory(userID.(string), memorialID, &req)
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
		Message: "创建成功",
		Data:    story,
	})
}

// GetLifeStories 获取生平故事列表
func (c *LifeStoryController) GetLifeStories(ctx *gin.Context) {
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

	stories, total, err := c.lifeStoryService.GetLifeStories(userID.(string), memorialID, page, pageSize)
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
			"list":      stories,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetLifeStory 获取生平故事详情
func (c *LifeStoryController) GetLifeStory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	storyID := ctx.Param("story_id")
	if storyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "故事ID不能为空",
		})
		return
	}

	story, err := c.lifeStoryService.GetLifeStory(userID.(string), storyID)
	if err != nil {
		if err.Error() == "生平故事不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
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
		Data:    story,
	})
}

// UpdateLifeStory 更新生平故事
func (c *LifeStoryController) UpdateLifeStory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	storyID := ctx.Param("story_id")
	if storyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "故事ID不能为空",
		})
		return
	}

	var req services.UpdateLifeStoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.lifeStoryService.UpdateLifeStory(userID.(string), storyID, &req)
	if err != nil {
		if err.Error() == "生平故事不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "无权修改此故事" {
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

// DeleteLifeStory 删除生平故事
func (c *LifeStoryController) DeleteLifeStory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	storyID := ctx.Param("story_id")
	if storyID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "故事ID不能为空",
		})
		return
	}

	err := c.lifeStoryService.DeleteLifeStory(userID.(string), storyID)
	if err != nil {
		if err.Error() == "生平故事不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "无权删除此故事" {
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

// CreateTimeline 创建时间轴事件
func (c *LifeStoryController) CreateTimeline(ctx *gin.Context) {
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

	var req services.CreateTimelineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	timeline, err := c.lifeStoryService.CreateTimeline(userID.(string), memorialID, &req)
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
		Message: "创建成功",
		Data:    timeline,
	})
}

// GetTimeline 获取时间轴
func (c *LifeStoryController) GetTimeline(ctx *gin.Context) {
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

	timelines, err := c.lifeStoryService.GetTimeline(userID.(string), memorialID)
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
		Data:    timelines,
	})
}

// DeleteTimeline 删除时间轴事件
func (c *LifeStoryController) DeleteTimeline(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	timelineID := ctx.Param("timeline_id")
	if timelineID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "时间轴事件ID不能为空",
		})
		return
	}

	err := c.lifeStoryService.DeleteTimeline(userID.(string), timelineID)
	if err != nil {
		if err.Error() == "时间轴事件不存在" {
			ctx.JSON(http.StatusNotFound, APIResponse{
				Code:    1004,
				Message: err.Error(),
			})
		} else if err.Error() == "无权删除此事件" {
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

// GetStoriesByCategory 按分类获取生平故事
func (c *LifeStoryController) GetStoriesByCategory(ctx *gin.Context) {
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

	category := ctx.Query("category")

	stories, err := c.lifeStoryService.GetStoriesByCategory(userID.(string), memorialID, category)
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
		Data:    stories,
	})
}