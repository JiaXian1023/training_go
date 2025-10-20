package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 票務模型
type Ticket struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Total     int    `gorm:"not null"`     // 總票數
	Remaining int    `gorm:"not null"`     // 剩餘票數
	Version   int    `gorm:"default:0"`    // 樂觀鎖版本號
}

// 訂單模型
type Order struct {
	gorm.Model
	UserID    uint
	TicketID  uint
	Quantity  int
	Status    string // "pending", "success", "failed"
}

var (
	db     *gorm.DB
	dbOnce sync.Once
)

// 初始化數據庫連接
func getDB() *gorm.DB {
	dbOnce.Do(func() {
		dsn := "username:password@tcp(127.0.0.1:3306)/ticket_db?charset=utf8mb4&parseTime=True&loc=Local"
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		// 自動遷移模型
		db.AutoMigrate(&Ticket{}, &Order{})
	})
	return db
}

func main() {
	// 初始化數據庫
	getDB()

	// 創建一個測試票務
	createTestTicket()

	r := gin.Default()

	// 搶票API
	r.POST("/api/ticket/:id/buy", buyTicket)

	// 啟動服務器
	r.Run(":8088")
}

// 創建測試票務
func createTestTicket() {
	db := getDB()
	var count int64
	db.Model(&Ticket{}).Where("name = ?", "Concert Ticket").Count(&count)
	if count == 0 {
		db.Create(&Ticket{
			Name:      "Concert Ticket",
			Total:     1000,
			Remaining: 1000,
		})
	}
}

// 搶票處理函數
func buyTicket(c *gin.Context) {
	// 獲取參數
	ticketID := c.Param("id")
	userID, _ := c.GetQuery("user_id")

	// 使用WaitGroup等待搶票結果
	var wg sync.WaitGroup
	wg.Add(1)

	// 使用channel傳遞搶票結果
	resultChan := make(chan string, 1)

	// 啟動goroutine處理搶票
	go func() {
		defer wg.Done()
		success, err := processTicketOrder(ticketID, userID)
		if err != nil {
			resultChan <- fmt.Sprintf("error: %v", err)
			return
		}
		if success {
			resultChan <- "success"
		} else {
			resultChan <- "failed: no tickets left"
		}
	}()

	// 設置超時
	select {
	case result := <-resultChan:
		if result == "success" {
			c.JSON(http.StatusOK, gin.H{"message": "搶票成功"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": result})
		}
	case <-time.After(5 * time.Second):
		c.JSON(http.StatusRequestTimeout, gin.H{"message": "請求超時"})
	}

	wg.Wait()
	close(resultChan)
}

// 處理訂單 - 使用樂觀鎖防止超賣
func processTicketOrder(ticketID, userID string) (bool, error) {
	db := getDB()

	// 開啟事務
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查詢票務信息
	var ticket Ticket
	if err := tx.Where("id = ?", ticketID).First(&ticket).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("票務不存在")
	}

	// 檢查庫存
	if ticket.Remaining <= 0 {
		tx.Rollback()
		return false, nil
	}

	// 更新庫存 - 使用樂觀鎖
	result := tx.Model(&Ticket{}).
		Where("id = ? AND version = ?", ticket.ID, ticket.Version).
		Updates(map[string]interface{}{
			"remaining": ticket.Remaining - 1,
			"version":   ticket.Version + 1,
		})

	if result.RowsAffected == 0 {
		tx.Rollback()
		return false, nil // 搶票失敗，可能是版本號不匹配
	}

	// 創建訂單
	order := Order{
		UserID:   parseUint(userID),
		TicketID: parseUint(ticketID),
		Quantity: 1,
		Status:   "success",
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		return false, err
	}

	return true, nil
}

// 輔助函數：字符串轉uint
func parseUint(s string) uint {
	var i uint
	fmt.Sscanf(s, "%d", &i)
	return i
}