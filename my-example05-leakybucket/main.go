package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// LeakyBucket 漏桶結構
type LeakyBucket struct {
	capacity  int           // 桶的容量
	remaining int           // 桶中剩餘的請求量
	rate      time.Duration // 漏水速率(處理請求的間隔)
	last      time.Time     // 上次處理請求的時間
	mu        sync.Mutex    // 互斥鎖
}

// NewLeakyBucket 創建一個新的漏桶
func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity:  capacity,
		remaining: capacity,
		rate:      rate,
		last:      time.Now(),
	}
}

// Allow 檢查是否允許通過請求
func (b *LeakyBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	
	// 計算從上次到現在漏出的水量
	elapsed := now.Sub(b.last)
	leaked := int(elapsed / b.rate)
	
	if leaked > 0 {
		b.remaining += leaked
		if b.remaining > b.capacity {
			b.remaining = b.capacity
		}
		b.last = now
	}
	
	if b.remaining > 0 {
		b.remaining--
		return true
	}
	
	return false
}

func main() {
	r := gin.Default()

	// 創建漏桶: 容量為10，速率為每秒處理2個請求(每500ms一個)
	bucket := NewLeakyBucket(10, 500*time.Millisecond)

	// 需要限流的API
	r.GET("/api", func(c *gin.Context) {
		if !bucket.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "too many requests",
				"status":  "rate limited",
			})
			return
		}

		// 正常處理請求
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    "your request has been processed",
		})
	})

	// 健康檢查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":8080")
}