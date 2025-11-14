package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileUploadManager 文件上传管理器
type FileUploadManager struct {
	uploadDir string
	maxSize   int64 // 最大文件大小（字节）
}

// NewFileUploadManager 创建文件上传管理器
func NewFileUploadManager(uploadDir string, maxSize int64) *FileUploadManager {
	return &FileUploadManager{
		uploadDir: uploadDir,
		maxSize:   maxSize,
	}
}

// UploadFile 上传文件
func (f *FileUploadManager) UploadFile(file *multipart.FileHeader, subDir string) (string, error) {
	// 检查文件大小
	if file.Size > f.maxSize {
		return "", fmt.Errorf("文件大小超过限制，最大允许 %d MB", f.maxSize/(1024*1024))
	}

	// 检查文件类型
	if !f.isAllowedFileType(file.Filename) {
		return "", fmt.Errorf("不支持的文件类型")
	}

	// 创建上传目录
	uploadPath := filepath.Join(f.uploadDir, subDir)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("创建上传目录失败: %v", err)
	}

	// 生成唯一文件名
	filename := f.generateUniqueFilename(file.Filename)
	fullPath := filepath.Join(uploadPath, filename)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开上传文件失败: %v", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	// 返回相对路径
	relativePath := filepath.Join(subDir, filename)
	return relativePath, nil
}

// DeleteFile 删除文件
func (f *FileUploadManager) DeleteFile(relativePath string) error {
	fullPath := filepath.Join(f.uploadDir, relativePath)
	if err := os.Remove(fullPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("删除文件失败: %v", err)
		}
	}
	return nil
}

// GetFileURL 获取文件访问URL
func (f *FileUploadManager) GetFileURL(relativePath string) string {
	// 这里应该返回完整的URL，包括域名
	// 在实际部署时需要配置正确的域名
	return "/uploads/" + strings.ReplaceAll(relativePath, "\\", "/")
}

// isAllowedFileType 检查是否为允许的文件类型
func (f *FileUploadManager) isAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedTypes := map[string]bool{
		// 图片类型
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		// 视频类型
		".mp4":  true,
		".avi":  true,
		".mov":  true,
		".wmv":  true,
		".flv":  true,
		".webm": true,
		// 音频类型
		".mp3":  true,
		".wav":  true,
		".aac":  true,
		".ogg":  true,
		".m4a":  true,
	}
	return allowedTypes[ext]
}

// generateUniqueFilename 生成唯一文件名
func (f *FileUploadManager) generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Format("20060102150405")
	uuid := GenerateUUID()[:8] // 使用UUID的前8位
	return fmt.Sprintf("%s_%s%s", timestamp, uuid, ext)
}

// GetFileType 根据文件扩展名获取文件类型
func GetFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	
	imageTypes := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	videoTypes := map[string]bool{
		".mp4": true, ".avi": true, ".mov": true, ".wmv": true, ".flv": true, ".webm": true,
	}
	audioTypes := map[string]bool{
		".mp3": true, ".wav": true, ".aac": true, ".ogg": true, ".m4a": true,
	}

	if imageTypes[ext] {
		return "image"
	} else if videoTypes[ext] {
		return "video"
	} else if audioTypes[ext] {
		return "audio"
	}
	
	return "unknown"
}

// ValidateImageFile 验证图片文件
func ValidateImageFile(file *multipart.FileHeader) error {
	// 检查文件大小（图片最大10MB）
	if file.Size > 10*1024*1024 {
		return fmt.Errorf("图片文件大小不能超过10MB")
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return nil
		}
	}
	
	return fmt.Errorf("只支持 JPG、PNG、GIF、WebP 格式的图片")
}

// ValidateVideoFile 验证视频文件
func ValidateVideoFile(file *multipart.FileHeader) error {
	// 检查文件大小（视频最大100MB）
	if file.Size > 100*1024*1024 {
		return fmt.Errorf("视频文件大小不能超过100MB")
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedTypes := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm"}
	
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return nil
		}
	}
	
	return fmt.Errorf("只支持 MP4、AVI、MOV、WMV、FLV、WebM 格式的视频")
}

// ValidateAudioFile 验证音频文件
func ValidateAudioFile(file *multipart.FileHeader) error {
	// 检查文件大小（音频最大50MB）
	if file.Size > 50*1024*1024 {
		return fmt.Errorf("音频文件大小不能超过50MB")
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedTypes := []string{".mp3", ".wav", ".aac", ".ogg", ".m4a"}
	
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return nil
		}
	}
	
	return fmt.Errorf("只支持 MP3、WAV、AAC、OGG、M4A 格式的音频")
}