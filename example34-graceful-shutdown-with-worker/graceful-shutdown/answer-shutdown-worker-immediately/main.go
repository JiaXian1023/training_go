package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"fmt"
)

// Consumer struct
type Consumer struct {
	inputChan chan int
	jobsChan  chan int
}

func getRandomTime() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10)
}

func withContextFunc(ctx context.Context, f func()) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(c)

		select {
		case <-ctx.Done():
			fmt.Println("withContextFunc-ctx.Done")
		case <-c:
			cancel()
			fmt.Println("stop")
			f()
		}
	}()

	return ctx
}

func (c *Consumer) queue(input int) bool {
	select {
	case c.inputChan <- input:
		log.Println("already send input value:", input)
		return true
	default:
		return false
	}
}
//監聽job執行
func (c Consumer) startConsumer(ctx context.Context) {
	for {
		select {
		case job := <-c.inputChan:
			c.jobsChan <- job
			if ctx.Err() != nil {
				fmt.Println("@startConsumer close")
				close(c.jobsChan)
				return
			}
		case <-ctx.Done():
			fmt.Println("@startConsumer done")
			close(c.jobsChan)
			return
		}
	}
}

//處理work
func (c *Consumer) process(num, job int) {
	n := getRandomTime()
	log.Printf("Sleeping %d seconds...\n", n)
	time.Sleep(time.Duration(n) * time.Second)
	log.Println("process worker:", num, " job value:", job)
}

//工作協程
func (c *Consumer) worker(ctx context.Context, num int, wg *sync.WaitGroup) {
	//執行
	defer wg.Done()
	log.Println("start the worker", num)
	for {
		select {
		case job := <-c.jobsChan:
			if ctx.Err() != nil {
				log.Println("get next job", job, "and close the worker", num)
				return
			}
			c.process(num, job)
		case <-ctx.Done():
			log.Println("close the worker", num)
			return
		}
	}
}

const poolSize = 2

func main() {
	finished := make(chan bool)
	wg := &sync.WaitGroup{}

	//wait 如果有兩個queue,就會等待2個執行
	wg.Add(poolSize)
	// create the consumer
	consumer := Consumer{
		inputChan: make(chan int, 10),
		jobsChan:  make(chan int, poolSize),
	}

	//關閉後才結束
	ctx := withContextFunc(context.Background(), func() {
		log.Println("cancel from ctrl+c event")
		wg.Wait()
		close(finished)
	})

	for i := 0; i < poolSize; i++ {
		go consumer.worker(ctx, i, wg)
	}

	go consumer.startConsumer(ctx)

	go func() {
		consumer.queue(1)
		consumer.queue(2)
		consumer.queue(3)
		consumer.queue(4)
	    consumer.queue(5)
		// consumer.queue(6)
		// consumer.queue(7)
		// consumer.queue(8)
		// consumer.queue(9)
		// consumer.queue(10)
	}()

	<-finished
	log.Println("Game over")
}
