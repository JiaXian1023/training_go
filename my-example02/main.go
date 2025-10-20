package main

import (
	"log"
	"net/http" 
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// TicketLimiter struct 包含一個 rate.Limiter，用於管理限流
type TicketLimiter struct {
	limiter *rate.Limiter
}

// NewTicketLimiter 創建一個新的 TicketLimiter 實例
// r: 每秒生成的令牌數，這裡設定為 100
// b: 桶的容量，這裡設定為 100
func NewTicketLimiter(r rate.Limit, b int) *TicketLimiter {
	return &TicketLimiter{
		limiter: rate.NewLimiter(r, b),
	}
}

// CheckLimit 是 Gin 的中介層 (middleware)，用於檢查是否超出了速率限制
func (tl *TicketLimiter) CheckLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !tl.limiter.Allow() {
			// 如果令牌桶中沒有令牌，則返回 429 Too Many Requests
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		// 如果有令牌，繼續處理請求
		c.Next()
	}
}

// getTicketHandler 處理領票請求
func getTicketHandler(c *gin.Context) {
	// 在這裡處理領票的業務邏輯，例如：
	// 檢查票券庫存、使用者驗證等
	log.Println("A ticket has been successfully claimed.")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Ticket claimed successfully!",
	})
}

func main() {
	// 創建一個 Gin 路由器
	r := gin.Default()

	// 創建一個令牌桶限流器，設定每秒生成 100 個令牌，桶容量為 100
	// 這表示在短時間內，最多允許 100 個請求通過
	limiter := NewTicketLimiter(rate.Limit(100), 100)

	// 使用中介層將限流器應用到 /get-ticket 路徑
	// 只有通過限流檢查的請求才能進入 getTicketHandler
	r.GET("/get-ticket", limiter.CheckLimit(), getTicketHandler)

	// 啟動伺服器
	log.Println("Server is running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}