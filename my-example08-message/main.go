package main

import (
	"fmt"
	"sync"
	"time" 
)

type TextMessage struct {
	ID      int
	Content string
}

func main() {
	// 创建消息通道
	messageChan := make(chan TextMessage, 100)

	// 启动消费者
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		fmt.Println("goroutine1")
		defer wg.Done()
		fmt.Println("goroutine11")
		for msg := range messageChan {
			processMessage(msg,"111")
		}
	}()

	go func() {
		fmt.Println("goroutine2")
		defer wg.Done()
		fmt.Println("goroutine22")
		for msg := range messageChan {
			processMessage(msg,"222")
		}
	}()

	// 生产消息
	for i := 1; i <= 10; i++ {
		fmt.Println("send ",i)
		messageChan <- TextMessage{
			ID:      i,
			Content: fmt.Sprintf("Message %d", i),
		}
	}

	close(messageChan)
	wg.Wait()
}

func processMessage(msg TextMessage,m string) {
 
	time.Sleep(500 * time.Millisecond) // 模拟处理时间
	fmt.Printf("processMessage:%s,Processed message %d: %s\n", m,msg.ID, msg.Content)
}