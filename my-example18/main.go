
package main

import (
    "fmt"
    "time"
)

// 演示两种主要差异
func main() {
    fmt.Println("=== 无缓冲 vs 缓冲 Channel ===")
    demoBuffering()
    
    // fmt.Println("\n=== 单向 Channel ===")
    // demoDirectional()
}

func demoBuffering() {
    // 无缓冲 channel
    unbuffered := make(chan string)
    
    go func() {
		fmt.Println("无缓冲: 接收前")
    	msg := <-unbuffered
		
    	fmt.Println("无缓冲: 接收到", msg)
        
    }()
    
   time.Sleep(1 * time.Second)
	fmt.Println("无缓冲: 发送前")
    unbuffered <- "Hello"
    fmt.Println("无缓冲: 发送后") // 这行会立即执行吗？
    
    // 缓冲 channel
    // buffered := make(chan string, 1)
    
    // go func() {
    //     fmt.Println("缓冲: 发送前")
    //     buffered <- "World"
    //     fmt.Println("缓冲: 发送后") // 这行会立即执行
    // }()
    
    // time.Sleep(1 * time.Second)
    // fmt.Println("缓冲: 接收前")
    // msg = <-buffered
    // fmt.Println("缓冲: 接收到", msg)
}

func demoDirectional() {
    ch := make(chan int, 2)
    
    // 只发送的函数
    sendOnly := func(ch chan<- int) {
        for i := 0; i < 3; i++ {
            ch <- i * 10
        }
        close(ch)
    }
    
    // 只接收的函数
    receiveOnly := func(ch <-chan int) {
        for num := range ch {
            fmt.Println("读取:", num)
        }
    }
    
    go sendOnly(ch)
    receiveOnly(ch)
}