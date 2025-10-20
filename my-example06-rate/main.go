package main

import (
	"net/http"
	"golang.org/x/time/rate"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiterMiddleware 創建限流中間件
func RateLimiterMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "too many requests",
				"status":  "rate limited",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()

	// 應用限流中間件: 每秒2個請求，突發容量10
	r.Use(RateLimiterMiddleware(2, 10))

	// 需要限流的API
	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    "your request has been processed",
		})
	})

	// 健康檢查接口(不受限流影響)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":8080")
}