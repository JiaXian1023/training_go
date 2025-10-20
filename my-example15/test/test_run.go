// test_integration.go
package main

import (
    "context"
    "fmt"
    "log"
)

func testIntegration() {
    sink, _ := NewRedisSinkConnector("localhost:9092", "user-updates", "test-group", "localhost:6379")
    ctx := context.Background()
    
    // 從 Redis 讀取用戶數據
    userData, err := sink.GetUserFromRedis(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("User data from Redis: %+v\n", userData)
}