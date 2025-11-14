package controllers

import (
	"net/http"
	"strconv"
	"yun-nian-memorial/internal/services"

	"github.com/gin-gonic/gin"
)

type AlbumController struct {
	albumService *services.AlbumService
}

func NewAlbumController(albumService *services.AlbumService) *AlbumController {
	return &AlbumController{
		albumService: albumService,
	}
}

// CreateAlbum 创建相册
func (c *AlbumController) CreateAlbum(ctx *gin.Context) {
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

	var req services.CreateAlbumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	album, err := c.albumService.CreateAlbum(userID.(string), memorialID, &req)
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
		Data:    album,
	})
}

// GetAlbums 获取相册列表
func (c *AlbumController) GetAlbums(ctx *gin.Context) {
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

	albums, total, err := c.albumService.GetAlbums(userID.(string), memorialID, page, pageSize)
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
			"list":      albums,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetAlbum 获取相册详情
func (c *AlbumController) GetAlbum(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	albumID := ctx.Param("album_id")
	if albumID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "相册ID不能为空",
		})
		return
	}

	album, err := c.albumService.GetAlbum(userID.(string), albumID)
	if err != nil {
		if err.Error() == "相册不存在" {
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
		Data:    album,
	})
}

// UpdateAlbum 更新相册
func (c *AlbumController) UpdateAlbum(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	albumID := ctx.Param("album_id")
	if albumID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "相册ID不能为空",
		})
		return
	}

	var req services.UpdateAlbumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.albumService.UpdateAlbum(userID.(string), albumID, &req)
	if err != nil {
		if err.Error() == "相册不存在" {
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
		Message: "更新成功",
	})
}

// DeleteAlbum 删除相册
func (c *AlbumController) DeleteAlbum(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	albumID := ctx.Param("album_id")
	if albumID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "相册ID不能为空",
		})
		return
	}

	err := c.albumService.DeleteAlbum(userID.(string), albumID)
	if err != nil {
		if err.Error() == "相册不存在" {
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
		Message: "删除成功",
	})
}

// AddPhoto 添加照片到相册
func (c *AlbumController) AddPhoto(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	albumID := ctx.Param("album_id")
	if albumID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "相册ID不能为空",
		})
		return
	}

	var req services.AddPhotoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	photo, err := c.albumService.AddPhoto(userID.(string), albumID, &req)
	if err != nil {
		if err.Error() == "相册不存在" {
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
		Message: "添加成功",
		Data:    photo,
	})
}

// UpdatePhoto 更新照片信息
func (c *AlbumController) UpdatePhoto(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	photoID := ctx.Param("photo_id")
	if photoID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "照片ID不能为空",
		})
		return
	}

	var req services.AddPhotoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	err := c.albumService.UpdatePhoto(userID.(string), photoID, &req)
	if err != nil {
		if err.Error() == "照片不存在" {
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
		Message: "更新成功",
	})
}

// DeletePhoto 删除照片
func (c *AlbumController) DeletePhoto(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIResponse{
			Code:    1002,
			Message: "用户未登录",
		})
		return
	}

	photoID := ctx.Param("photo_id")
	if photoID == "" {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    1001,
			Message: "照片ID不能为空",
		})
		return
	}

	err := c.albumService.DeletePhoto(userID.(string), photoID)
	if err != nil {
		if err.Error() == "照片不存在" {
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
		Message: "删除成功",
	})
}