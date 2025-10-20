package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Consumer struct
//設計狀態
type Consumer struct {
	inputChan chan int
	jobsChan  chan int
}

func getRandomTime() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(5)
}
//發送
func (c *Consumer) queue(input int) bool {
	c.jobsChan <- input
	fmt.Println("already send input value:", input)
	return true

	// select {
	// case c.jobsChan <- input:
	// 	fmt.Println("already send input value:", input)
	// 	return true
	// default:
	// 	return false
	// }
}
//接收
func (c *Consumer) worker(num int) {
	for job := range c.jobsChan {
		n := getRandomTime()
		fmt.Printf("@worker:%v @Sleeping %d seconds...\n", num,  n)
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Println("@worker:", num, " job value:", job)
	}
}

const poolSize = 2
//worker接收jobsChan  goroutine xn等待
//cobsumer 1 ,2 發送
//worker其中的1,2工作攜程會接收到數據, 工作完畢後結束主程式
func main() {
	// create the consumer
	consumer := Consumer{
		inputChan: make(chan int, 1),
		jobsChan:  make(chan int, poolSize),
	}
	//只有兩個worker,0,1
	for i := 0; i < poolSize; i++ {
		fmt.Println("i",i)
		go consumer.worker(i)
	}

	consumer.queue(1) 
	consumer.queue(2) 
	consumer.queue(3) 
	time.Sleep(5 * time.Second)
}
