package services

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
	"yun-nian-memorial/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MonitorService struct {
	db *gorm.DB
}

func NewMonitorService(db *gorm.DB) *MonitorService {
	return &MonitorService{
		db: db,
	}
}

// MetricType 指标类型
const (
	MetricTypeCPU      = "cpu"
	MetricTypeMemory   = "memory"
	MetricTypeDisk     = "disk"
	MetricTypeAPI      = "api"
	MetricTypeDatabase = "database"
)

// RecordMetric 记录监控指标
func (s *MonitorService) RecordMetric(metricType string, value float64, unit string, additionalInfo map[string]interface{}) error {
	infoJSON, _ := json.Marshal(additionalInfo)
	
	metric := &models.SystemMonitor{
		ID:             uuid.New().String(),
		MetricType:     metricType,
		MetricValue:    value,
		MetricUnit:     unit,
		AdditionalInfo: string(infoJSON),
		CreatedAt:      time.Now(),
	}
	
	return s.db.Create(metric).Error
}

// RecordSystemMetrics 记录系统指标
func (s *MonitorService) RecordSystemMetrics() error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// 记录内存使用
	memoryUsageMB := float64(m.Alloc) / 1024 / 1024
	if err := s.RecordMetric(MetricTypeMemory, memoryUsageMB, "MB", map[string]interface{}{
		"total_alloc": m.TotalAlloc,
		"sys":         m.Sys,
		"num_gc":      m.NumGC,
	}); err != nil {
		return err
	}
	
	// 记录CPU使用（Goroutine数量作为简单指标）
	numGoroutines := float64(runtime.NumGoroutine())
	if err := s.RecordMetric(MetricTypeCPU, numGoroutines, "goroutines", map[string]interface{}{
		"num_cpu": runtime.NumCPU(),
	}); err != nil {
		return err
	}
	
	return nil
}

// RecordAPIMetrics 记录API指标
func (s *MonitorService) RecordAPIMetrics(endpoint string, method string, statusCode int, duration int64) error {
	return s.RecordMetric(MetricTypeAPI, float64(duration), "ms", map[string]interface{}{
		"endpoint":    endpoint,
		"method":      method,
		"status_code": statusCode,
	})
}

// RecordDatabaseMetrics 记录数据库指标
func (s *MonitorService) RecordDatabaseMetrics() error {
	// 获取数据库连接池状态
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	
	stats := sqlDB.Stats()
	
	// 记录打开的连接数
	if err := s.RecordMetric(MetricTypeDatabase, float64(stats.OpenConnections), "connections", map[string]interface{}{
		"in_use":      stats.InUse,
		"idle":        stats.Idle,
		"wait_count":  stats.WaitCount,
		"wait_duration": stats.WaitDuration.Milliseconds(),
	}); err != nil {
		return err
	}
	
	return nil
}

// GetMetrics 获取监控指标
func (s *MonitorService) GetMetrics(metricType string, startTime, endTime time.Time, limit int) ([]models.SystemMonitor, error) {
	var metrics []models.SystemMonitor
	
	query := s.db.Model(&models.SystemMonitor{})
	
	if metricType != "" {
		query = query.Where("metric_type = ?", metricType)
	}
	
	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Order("created_at DESC").Find(&metrics).Error
	return metrics, err
}

// GetMetricStats 获取指标统计
func (s *MonitorService) GetMetricStats(metricType string, duration time.Duration) (map[string]interface{}, error) {
	startTime := time.Now().Add(-duration)
	
	var metrics []models.SystemMonitor
	err := s.db.Where("metric_type = ? AND created_at >= ?", metricType, startTime).
		Order("created_at ASC").
		Find(&metrics).Error
	if err != nil {
		return nil, err
	}
	
	if len(metrics) == 0 {
		return map[string]interface{}{
			"count":   0,
			"average": 0,
			"min":     0,
			"max":     0,
		}, nil
	}
	
	// 计算统计数据
	var sum, min, max float64
	min = metrics[0].MetricValue
	max = metrics[0].MetricValue
	
	for _, metric := range metrics {
		sum += metric.MetricValue
		if metric.MetricValue < min {
			min = metric.MetricValue
		}
		if metric.MetricValue > max {
			max = metric.MetricValue
		}
	}
	
	average := sum / float64(len(metrics))
	
	return map[string]interface{}{
		"count":      len(metrics),
		"average":    average,
		"min":        min,
		"max":        max,
		"unit":       metrics[0].MetricUnit,
		"start_time": startTime,
		"end_time":   time.Now(),
	}, nil
}

// GetSystemHealth 获取系统健康状态
func (s *MonitorService) GetSystemHealth() (map[string]interface{}, error) {
	health := make(map[string]interface{})
	
	// 检查数据库连接
	sqlDB, err := s.db.DB()
	if err != nil {
		health["database"] = map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	} else {
		if err := sqlDB.Ping(); err != nil {
			health["database"] = map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
		} else {
			stats := sqlDB.Stats()
			health["database"] = map[string]interface{}{
				"status":           "healthy",
				"open_connections": stats.OpenConnections,
				"in_use":           stats.InUse,
				"idle":             stats.Idle,
			}
		}
	}
	
	// 获取内存使用情况
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	health["memory"] = map[string]interface{}{
		"status":      "healthy",
		"alloc_mb":    float64(m.Alloc) / 1024 / 1024,
		"total_alloc": float64(m.TotalAlloc) / 1024 / 1024,
		"sys_mb":      float64(m.Sys) / 1024 / 1024,
		"num_gc":      m.NumGC,
	}
	
	// 获取CPU信息
	health["cpu"] = map[string]interface{}{
		"status":        "healthy",
		"num_cpu":       runtime.NumCPU(),
		"num_goroutine": runtime.NumGoroutine(),
	}
	
	// 获取最近的API性能
	recentAPIMetrics, _ := s.GetMetrics(MetricTypeAPI, time.Now().Add(-5*time.Minute), time.Time{}, 100)
	if len(recentAPIMetrics) > 0 {
		var totalDuration float64
		for _, metric := range recentAPIMetrics {
			totalDuration += metric.MetricValue
		}
		avgDuration := totalDuration / float64(len(recentAPIMetrics))
		
		health["api"] = map[string]interface{}{
			"status":           "healthy",
			"avg_response_ms":  avgDuration,
			"recent_requests":  len(recentAPIMetrics),
		}
	} else {
		health["api"] = map[string]interface{}{
			"status": "no_data",
		}
	}
	
	// 整体健康状态
	health["overall_status"] = "healthy"
	health["timestamp"] = time.Now()
	
	return health, nil
}

// GetDashboardStats 获取仪表板统计数据
func (s *MonitorService) GetDashboardStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// 获取最近1小时的系统指标统计
	memoryStats, _ := s.GetMetricStats(MetricTypeMemory, 1*time.Hour)
	stats["memory"] = memoryStats
	
	cpuStats, _ := s.GetMetricStats(MetricTypeCPU, 1*time.Hour)
	stats["cpu"] = cpuStats
	
	apiStats, _ := s.GetMetricStats(MetricTypeAPI, 1*time.Hour)
	stats["api"] = apiStats
	
	dbStats, _ := s.GetMetricStats(MetricTypeDatabase, 1*time.Hour)
	stats["database"] = dbStats
	
	// 获取系统健康状态
	health, _ := s.GetSystemHealth()
	stats["health"] = health
	
	return stats, nil
}

// CleanOldMetrics 清理旧的监控数据
func (s *MonitorService) CleanOldMetrics(daysToKeep int) error {
	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)
	
	result := s.db.Where("created_at < ?", cutoffDate).Delete(&models.SystemMonitor{})
	if result.Error != nil {
		return fmt.Errorf("清理旧监控数据失败: %v", result.Error)
	}
	
	fmt.Printf("已清理 %d 条旧监控数据\n", result.RowsAffected)
	return nil
}

// GetMetricTrend 获取指标趋势
func (s *MonitorService) GetMetricTrend(metricType string, duration time.Duration, interval time.Duration) ([]map[string]interface{}, error) {
	startTime := time.Now().Add(-duration)
	
	var metrics []models.SystemMonitor
	err := s.db.Where("metric_type = ? AND created_at >= ?", metricType, startTime).
		Order("created_at ASC").
		Find(&metrics).Error
	if err != nil {
		return nil, err
	}
	
	// 按时间间隔分组计算平均值
	trendData := make([]map[string]interface{}, 0)
	
	if len(metrics) == 0 {
		return trendData, nil
	}
	
	currentBucket := startTime
	var bucketValues []float64
	
	for _, metric := range metrics {
		// 如果超过当前时间桶，计算平均值并开始新桶
		if metric.CreatedAt.After(currentBucket.Add(interval)) {
			if len(bucketValues) > 0 {
				var sum float64
				for _, v := range bucketValues {
					sum += v
				}
				avg := sum / float64(len(bucketValues))
				
				trendData = append(trendData, map[string]interface{}{
					"timestamp": currentBucket,
					"value":     avg,
					"count":     len(bucketValues),
				})
			}
			
			// 移动到下一个时间桶
			currentBucket = currentBucket.Add(interval)
			bucketValues = []float64{}
		}
		
		bucketValues = append(bucketValues, metric.MetricValue)
	}
	
	// 处理最后一个桶
	if len(bucketValues) > 0 {
		var sum float64
		for _, v := range bucketValues {
			sum += v
		}
		avg := sum / float64(len(bucketValues))
		
		trendData = append(trendData, map[string]interface{}{
			"timestamp": currentBucket,
			"value":     avg,
			"count":     len(bucketValues),
		})
	}
	
	return trendData, nil
}

// StartMonitoring 启动定期监控
func (s *MonitorService) StartMonitoring(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			// 记录系统指标
			if err := s.RecordSystemMetrics(); err != nil {
				fmt.Printf("记录系统指标失败: %v\n", err)
			}
			
			// 记录数据库指标
			if err := s.RecordDatabaseMetrics(); err != nil {
				fmt.Printf("记录数据库指标失败: %v\n", err)
			}
		}
	}()
}
