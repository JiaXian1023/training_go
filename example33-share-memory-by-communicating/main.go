package main

import (
	"fmt"
	"sync"
)

func addByShareMemory(n int) []int {
	var ints []int
	var wg sync.WaitGroup
	var mux sync.Mutex

	wg.Add(n)
	fmt.Println("set")
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			mux.Lock()
			ints = append(ints, i)
			mux.Unlock()
		}(i)
	}
fmt.Println("wait")
	wg.Wait()

	return ints
}

func syncTest(n int )[] int{
	var ints [] int
	var wg sync.WaitGroup
	var mux sync.Mutex
	wg.Add(n)
	for i:=0 ;i<n ; i++{
		go func(i int){
			defer wg.Done()
			mux.Lock()
			ints = append(ints, i)
			mux.Unlock()
		}(i)
	}
	wg.Wait()
	return ints
}

func syncRun(n int)[] int{
	var ints [] int
	var wg sync.WaitGroup
	var mux sync.Mutex
	wg.Add(n)
	for i:=0 ; i<n;i++{
		go func (i int){
			defer wg.Done()
			mux.Lock()
			ints=append(ints,i)
			mux.Unlock()
		}(i)
	}
	wg.Wait()
	return ints
}

func addByShareCommunicate(n int) []int {
	var ints []int
	channel := make(chan int, n)

	for i := 0; i < n; i++ {
		fmt.Println("i",i)
		go func(channel chan<- int, order int) {
			channel <- order
		}(channel, i)
	}

	for i := range channel {
		fmt.Println("range",i)
		ints = append(ints, i)

		if len(ints) == n {
			break
		}
	}

	close(channel)

	return ints
}

func chanTest (n int )[] int{
	var ints []int
	channel :=make(chan int,n)

	for i:=0;i<n;i++{
		go func(channel chan<-int,order int){
			channel<-order
		}(channel,i)
	}

	for i:= range channel{
		ints=append(ints,i)
		if(len(ints)== n){
			break
		}
	}
	close (channel)
	return ints
}

func chanRun(n int)[] int{
	var ints []int
	channel:=make(chan int,n)
	for i:=0; i<n; i++{
		go func(channel chan<-int, order int){
			channel<-order

		}(channel,i)
	}
	for i:=range channel{
		ints=append(ints,i)
		if(len(ints)==n){
			break
		}
	}
	close(channel)
	return ints
}

func main() { 
	// foo := syncRun(10)
	// fmt.Println(len(foo))
	// fmt.Println(foo)

	foo := chanTest(10)
	fmt.Println(len(foo))
	fmt.Println(foo)
}
