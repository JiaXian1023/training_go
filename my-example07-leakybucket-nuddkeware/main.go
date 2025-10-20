package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// LeakyBucket 漏桶結構體
type LeakyBucket struct {
	capacity     int           // 桶的總容量
	remaining    int           // 當前剩餘容量
	rate         time.Duration // 漏水速率(處理間隔)
	lastLeakTime time.Time     // 上次漏水時間
	mu           sync.Mutex    // 互斥鎖保證併發安全
}

// NewLeakyBucket 創建漏桶實例
// capacity: 桶容量
// rate: 處理速率(每個請求間隔時間)
func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity:     capacity,
		remaining:    capacity,
		rate:         rate,
		lastLeakTime: time.Now(),
	}
}

// Allow 檢查是否允許請求通過
func (b *LeakyBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	
	// 計算從上次檢查到現在漏出的水量
	elapsed := now.Sub(b.lastLeakTime)
	leaked := int(elapsed / b.rate)
	
	// 更新漏桶狀態
	if leaked > 0 {
		b.remaining += leaked
		if b.remaining > b.capacity {
			b.remaining = b.capacity
		}
		b.lastLeakTime = now
	}
	
	// 判斷是否有足夠容量處理當前請求
	if b.remaining > 0 {
		b.remaining--
		return true
	}
	
	return false
}

// LeakyBucketMiddleware 漏桶限流中間件
func LeakyBucketMiddleware(capacity int, rate time.Duration) gin.HandlerFunc {
	bucket := NewLeakyBucket(capacity, rate)

	return func(c *gin.Context) {
		if !bucket.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "請求過於頻繁，請稍後再試",
				"data":    nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()

	// 應用漏桶限流中間件
	// 參數說明:
	// 10 - 桶容量(突發請求最大數量)
	// 100*time.Millisecond - 每100ms處理一個請求(即每秒10個)
	apiGroup := r.Group("/api").Use(LeakyBucketMiddleware(10, 100*time.Millisecond))
	{
		apiGroup.GET("/resource", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": "請求成功",
				"data":    "這是受保護的資源",
			})
		})
	}

	// 不受限流的健康檢查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "服務正常運行",
		})
	})

	// 啟動服務
	r.Run(":8080")
}


/*

//全域
func main() {
	r := gin.Default()
	
	// 全局中間件 - 對所有路由生效
	r.Use(LeakyBucketMiddleware(100, 10*time.Millisecond)) // 每秒100個請求
	
	// ... 路由定義
}

//路由組
func main() {
	r := gin.Default()
	
	// API路由組 - 特定限流策略
	api := r.Group("/api")
	api.Use(LeakyBucketMiddleware(50, 20*time.Millisecond)) // 每秒50個請求
	{
		api.GET("/users", getUserList)
		api.POST("/users", createUser)
	}
	
	// 公開路由 - 不限流
	r.GET("/public", getPublicInfo)
}


//單個
func main() {
	r := gin.Default()
	
	// 對特定路由應用限流
	r.GET("/limited",
		LeakyBucketMiddleware(1, time.Second), // 每秒1個請求
		func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "這是嚴格限流的接口"})
		})
}

*/