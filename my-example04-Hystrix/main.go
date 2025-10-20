package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 Hystrix 配置
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		Timeout:                1000, // 超時時間(毫秒)
		MaxConcurrentRequests:  10,   // 最大併發量
		RequestVolumeThreshold: 5,    // 觸發熔斷的最小請求數
		SleepWindow:            5000, // 熔斷後多久嘗試恢復(毫秒)
		ErrorPercentThreshold:  20,   // 錯誤百分比閾值
	})

	// 創建 Gin 路由
	r := gin.Default()

	// 健康檢查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 需要保護的 API
	r.GET("/api", func(c *gin.Context) {
		// 使用 Hystrix 保護這個 API
		err := hystrix.Do("my_command", func() error {
			// 這裡是業務邏輯
			if rand.Intn(10) < 3 { // 模擬30%的失敗率
				return fmt.Errorf("random error occurred")
			}

			// 模擬耗時操作
			time.Sleep(time.Duration(rand.Intn(800)) * time.Millisecond)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
			return nil
		}, func(err error) error {
			// 熔斷或錯誤時的降級處理
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "service unavailable (fallback)",
				"error":   err.Error(),
			})
			return nil
		})

		if err != nil {
			// 這裡的錯誤已經在 fallback 函數中處理
			return
		}
	})

	// 啟動 Hystrix 的指標流 (可選)
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(":8081", hystrixStreamHandler)

	// 啟動 Gin 服務
	r.Run(":8080")
}