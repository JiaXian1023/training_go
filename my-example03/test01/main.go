package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// 定義類型和常量
type QueryResult struct {
	Data  interface{} `json:"data"`
	Param string      `json:"param"`
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// 全局變數（實際項目中應該使用依賴注入）
var (
	// 緩存相關
	cache    sync.Map
	cacheTTL = 10 * time.Minute

	// 並發控制
	queryMutex sync.Mutex
	inProgress = make(map[string][]chan interface{})

	// 指標監控
	metrics = struct {
		queryDuration *prometheus.HistogramVec
		cacheHits     prometheus.Counter
		cacheMisses   prometheus.Counter
		queryErrors   *prometheus.CounterVec
		timeoutErrors prometheus.Counter
	}{
		queryDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "query_duration_seconds",
			Help:    "Time spent executing database queries",
			Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0},
		}, []string{"status"}),
		cacheHits: promauto.NewCounter(prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		}),
		cacheMisses: promauto.NewCounter(prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		}),
		queryErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "query_errors_total",
			Help: "Total number of query errors",
		}, []string{"type"}),
		timeoutErrors: promauto.NewCounter(prometheus.CounterOpts{
			Name: "timeout_errors_total",
			Help: "Total number of timeout errors",
		}),
	}
)

// 主處理函數
func asyncQueryHandler(c *gin.Context) {
	startTime := time.Now()

	// 參數驗證
	param := c.Query("param")
	if param == "" {
		respondWithError(c, http.StatusBadRequest, "param query parameter is required")
		return
	}

	if len(param) > 100 {
		respondWithError(c, http.StatusBadRequest, "param too long")
		return
	}

	// 1. 檢查緩存
	if result, found := getFromCache(param); found {
		metrics.cacheHits.Inc()
		log.Printf("Cache hit for param: %s", param)
		respondWithSuccess(c, result)
		return
	}
	metrics.cacheMisses.Inc()

	// 2. 檢查是否已有相同查詢在進行中
	if result, found := getFromInProgress(param); found {
		log.Printf("Reusing in-progress query for param: %s", param)
		respondWithSuccess(c, result)
		return
	}

	// 3. 執行異步查詢
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	result, err := executeQueryWithTimeout(ctx, param)
	if err != nil {
		handleQueryError(c, err, startTime)
		return
	}

	// 4. 設置緩存並響應
	setToCache(param, result)
	metrics.queryDuration.WithLabelValues("success").Observe(time.Since(startTime).Seconds())
	respondWithSuccess(c, result)
}

// 執行查詢（含超時控制）
func executeQueryWithTimeout(ctx context.Context, param string) (interface{}, error) {
	resultChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)

	// 註冊正在進行的查詢
	registerInProgress(param, resultChan)
	defer unregisterInProgress(param, resultChan)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Query panic recovered: %v", r)
				errChan <- fmt.Errorf("query panic: %v", r)
			}
		}()

		result, err := queryDatabase(ctx, param)
		if err != nil {
			errChan <- err
			return
		}

		select {
		case resultChan <- result:
			// 成功發送
		case <-ctx.Done():
			// 上下文已取消
		}
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// 數據庫查詢實現
func queryDatabase(ctx context.Context, param string) (interface{}, error) {
	// 模擬數據庫查詢延遲
	select {
	case <-time.After(time.Duration(100+len(param)%100) * time.Millisecond):
		// 模擬成功查詢
		return &QueryResult{
			Data:  fmt.Sprintf("Processed: %s", param),
			Param: param,
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// 緩存管理
func getFromCache(key string) (interface{}, bool) {
	if item, found := cache.Load(key); found {
		cacheItem := item.(cacheItem)
		if time.Now().Before(cacheItem.expiration) {
			return cacheItem.value, true
		}
		// 緩存過期，刪除
		cache.Delete(key)
	}
	return nil, false
}

func setToCache(key string, value interface{}) {
	item := cacheItem{
		value:      value,
		expiration: time.Now().Add(cacheTTL),
	}
	cache.Store(key, item)
}

// 並發查詢控制（防止緩存擊穿）
func getFromInProgress(param string) (interface{}, bool) {
	queryMutex.Lock()
	defer queryMutex.Unlock()

	if chans, exists := inProgress[param]; exists && len(chans) > 0 {
		// 創建新的 channel 來接收結果
		resultChan := make(chan interface{}, 1)
		inProgress[param] = append(inProgress[param], resultChan)
		return <-resultChan, true
	}
	return nil, false
}

func registerInProgress(param string, resultChan chan interface{}) {
	queryMutex.Lock()
	defer queryMutex.Unlock()

	inProgress[param] = append(inProgress[param], resultChan)
}

func unregisterInProgress(param string, resultChan chan interface{}) {
	queryMutex.Lock()
	defer queryMutex.Unlock()

	if chans, exists := inProgress[param]; exists {
		// 移除指定的 channel
		for i, ch := range chans {
			if ch == resultChan {
				inProgress[param] = append(chans[:i], chans[i+1:]...)
				break
			}
		}

		// 如果沒有等待的 channel，刪除該 key
		if len(inProgress[param]) == 0 {
			delete(inProgress, param)
		}
	}
}

// 通知所有等待相同查詢的請求
func notifyInProgress(param string, result interface{}) {
	queryMutex.Lock()
	defer queryMutex.Unlock()

	if chans, exists := inProgress[param]; exists {
		for _, ch := range chans {
			select {
			case ch <- result:
				// 成功發送
			default:
				// channel 已滿，跳過
			}
		}
		delete(inProgress, param)
	}
}

// 響應輔助函數
func respondWithSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      data,
		"timestamp": time.Now().Unix(),
	})
}

func respondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success":   false,
		"error":     message,
		"timestamp": time.Now().Unix(),
	})
}

func handleQueryError(c *gin.Context, err error, startTime time.Time) {
	log.Printf("Query error: %v", err)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		metrics.timeoutErrors.Inc()
		metrics.queryDuration.WithLabelValues("timeout").Observe(time.Since(startTime).Seconds())
		respondWithError(c, http.StatusGatewayTimeout, "request timeout")

	case errors.Is(err, context.Canceled):
		metrics.queryErrors.WithLabelValues("canceled").Inc()
		respondWithError(c, http.StatusServiceUnavailable, "request canceled")

	default:
		metrics.queryErrors.WithLabelValues("internal").Inc()
		metrics.queryDuration.WithLabelValues("error").Observe(time.Since(startTime).Seconds())
		respondWithError(c, http.StatusInternalServerError, "internal server error")
	}
}

// 健康檢查和中間件
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	})
}

func metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		// 記錄請求指標
		requestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			fmt.Sprintf("%d", status),
		).Observe(duration)
	}
}

var requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_duration_seconds",
	Help:    "HTTP request duration in seconds",
	Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0},
}, []string{"method", "path", "status"})

// 初始化函數
func init() {
	// 定期清理過期緩存
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			cleanupExpiredCache()
		}
	}()
}

func cleanupExpiredCache() {
	now := time.Now()
	cache.Range(func(key, value interface{}) bool {
		item := value.(cacheItem)
		if now.After(item.expiration) {
			cache.Delete(key)
		}
		return true
	})
}
