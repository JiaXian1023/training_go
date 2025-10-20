package main

import (
	"fmt"
	
)

func main() {


	ch := make(chan int)
	go func(){	
		ch <-1
	}()


	fmt.Println(<-ch)
 
/*
	ch := make(chan int, 2)
	go func() {
		
		fmt.Println("11")
		ch <- 1
		//ch <- 2
		// close(ch)

	}()

		go func() { 
		//fmt.Println("11")
		//ch <- 1
		fmt.Println("22")
		ch <- 2
		
	}()
	// fmt.Println("22")
	// for n := range ch {
	// 	fmt.Println(n)
	// }

	go func() {
		
		for n := range ch {
			fmt.Println(n)
		}
	}()
	defer  close(ch)

	time.Sleep(2 * time.Second)*/
}
