package main

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"sync"
	"time"
	"fmt"
)

var (
	db *sql.DB
	c  *cache.Cache
)

func init() {
	// 初始化數據庫連接池
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname?parseTime=true")
	if err != nil {
		panic(err)
	}
	
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// 初始化緩存
	c = cache.New(5*time.Minute, 10*time.Minute)
}

type rateLimiter struct {
	mu         sync.Mutex
	rate       int           // 每秒允許的請求數
	capacity   int           // 桶容量
	tokens     int           // 當前令牌數
	lastRefill time.Time     // 上次補充時間
}

func newRateLimiter(rate, capacity int) *rateLimiter {
	return &rateLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

func (rl *rateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	
	// 計算應該補充的令牌數
	tokensToAdd := int(elapsed.Seconds() * float64(rl.rate))
	if tokensToAdd > 0 {
		rl.tokens = min(rl.tokens+tokensToAdd, rl.capacity)
		rl.lastRefill = now
	}
	
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func rateLimiterMiddleware() gin.HandlerFunc {
	limiter := newRateLimiter(1000, 1000) // 每秒1000請求
	
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}
		c.Next()
	}
}

func cacheControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "public, max-age=60")
		c.Next()
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	
	r := gin.New()
	
	// 中間件
	r.Use(gin.Recovery())
	r.Use(rateLimiterMiddleware())
	r.Use(cacheControlMiddleware())
	
	// 路由
	r.GET("/api/query", queryHandler)
	r.GET("/api/async-query", asyncQueryHandler)
	
	// 優化HTTP服務器配置
	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("RUN",s.Addr)
	// 啟動服務器
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}

func queryHandler(c *gin.Context) {
	param := c.Query("param")
	
	// 檢查緩存
	if data, found := c.Get(param); found {
		c.JSON(http.StatusOK, data)
		return
	}
	
	// 查詢數據庫
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	result, err := queryDatabase(ctx, param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// 設置緩存 - 使用全局緩存對象而不是gin.Context
	c.Set(param, result)
	
	c.JSON(http.StatusOK, result)
}

//異步處理
func asyncQueryHandler(c *gin.Context) {
	param := c.Query("param")
	
	// 檢查緩存
	if data, found := c.Get(param); found {
		c.JSON(http.StatusOK, data)
		return
	}
	
	resultChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)
	
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		
		result, err := queryDatabase(ctx, param)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()
	
	select {
	case result := <-resultChan:
		c.Set(param, result)
		c.JSON(http.StatusOK, result)
	case err := <-errChan:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case <-time.After(5 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
	}
}

func queryDatabase(ctx context.Context, param string) (interface{}, error) {
	// 實際數據庫查詢邏輯
	var result string
	err := db.QueryRowContext(ctx, "SELECT data FROM table WHERE param = ?", param).Scan(&result)
	if err != nil {
		return nil, err
	}
	
	return gin.H{
		"param": param,
		"data":  result,
		"time":  time.Now().Unix(),
	}, nil
}