package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	outChan := make(chan int)
	errChan := make(chan error)
	finishChan := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(outChan chan<- int, errChan chan<- error, val int, wg *sync.WaitGroup) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			fmt.Println("finished job id:", val)
			outChan <- val
			fmt.Println("finished job in:", val)
			//outChan進出
			if val == 160 {
				errChan <- errors.New("error in 60")
			}

		}(outChan, errChan, i, &wg)
	}

	go func() {
		//fmt.Println("Wait")
		wg.Wait() //wg=0才會往下
		//fmt.Println("@wait")
		close(finishChan)
	}()

Loop:
	for {
		select {
		case val := <-outChan:
			//outChan出
			fmt.Println("finished out:", val)
		case err := <-errChan:
			fmt.Println("error:", err)
			break Loop
		case <-finishChan:
			fmt.Println("close")
			break Loop
		case <-time.After(100000 * time.Millisecond):
			fmt.Println("time out")
			break Loop
		}
	}
}
