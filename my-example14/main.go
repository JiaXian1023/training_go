package main

import (
	"fmt"   
	 "os"
    "os/signal"
    "syscall"

    "time"
)

// 12AB34CD56EF78GH910IJ1112KL1314MN1516OP1718QR1920ST2122UV2324WX2526YZ
//
// 2728
func main() {
	// number := make(chan bool)
	// letter := make(chan bool)
	// done := make(chan bool)

	// go func() {
	// 	i := 1
	// 	for {
	// 		select {
	// 		case <-number:
	// 			fmt.Print(i)
	// 			i++
	// 			fmt.Print(i)
	// 			i++
	// 			letter <- true
	// 		}
	// 	}
	// }()

	// go func() {
	// 	j := 'A'
	// 	for {
	// 		select {
	// 		case <-letter:
	// 			if j >= 'Z' {
	// 				done <- true
	// 			} else {
	// 				fmt.Print(string(j))
	// 				j++
	// 				fmt.Print(string(j))
	// 				j++
	// 				number <- true
	// 			}
	// 		}
	// 	}
	// }()

	// number <- true

	// for {
	// 	select {
	// 	case <-done:
	// 		//fmt.Println("wait")
	// 		return
	// 	}
	// }

	// letter, number := make(chan bool), make(chan bool)
	// wait := sync.WaitGroup{}

	// go func() {
	// 	i := 1
	// 	for {
	// 		select {
	// 		case <-number:
	// 			fmt.Print(i)
	// 			i++
	// 			fmt.Print(i)
	// 			i++
	// 			letter <- true
	// 		}
	// 	}
	// }()
	// wait.Add(1)
	// go func(wait *sync.WaitGroup) {
	// 	i := 'A'
	// 	for {
	// 		select {
	// 		case <-letter:
	// 			if i >= 'Z' {
	// 				wait.Done()
	// 				return
	// 			}

	// 			fmt.Print(string(i))
	// 			i++
	// 			fmt.Print(string(i))
	// 			i++
	// 			number <- true
	// 		}

	// 	}
	// }(&wait)
	// number <- true
	// wait.Wait() 
	
	
	// go func() {
	// 	i:=0
	// 	//fmt.Print("start1") 
	// 	for {
	// 		i++
	// 		fmt.Print("@",i) 
	// 		time.Sleep(30)
	// 	}
	// }()
	// fmt.Print("main")
	// time.Sleep(50000)


	// stopChan := make(chan string)
	// go func() {
    //     // 方式1：赋值给变量 - 会分配内存给 stop 变量
    //     select {
    //     case stop := <-stopChan:
    //         fmt.Printf("Value: %s\n", stop) // 使用接收到的值
    //     }
    // }()
	// go func() {
    //     // 方式2：不赋值 - 更轻量，不分配额外变量
    //     select {
    //     case <-stopChan:
    //         fmt.Println("Signal received") // 不关心具体值
    //     }
    // }()
    // stopChan <- "stop message"
  	// close(stopChan)  


	// fmt.Println("以下是数值的chan")
    // ci:=make(chan int,3)
    // ci<-1
    // close(ci) 
    // num,ok := <- ci
    // fmt.Printf("读chan的协程结束，num=%v， ok=%v\n",num,ok)
    // num1,ok1 := <-ci
    // fmt.Printf("再读chan的协程结束，num=%v， ok=%v\n",num1,ok1)
    // num2,ok2 := <-ci
    // fmt.Printf("再再读chan的协程结束，num=%v， ok=%v\n",num2,ok2)
    
    // fmt.Println("以下是字符串chan")
    // cs := make(chan string,3)
    // cs <- "aaa"
    // close(cs)
    // str,ok := <- cs
    // fmt.Printf("读chan的协程结束，str=%v， ok=%v\n",str,ok)
    // str1,ok1 := <-cs
    // fmt.Printf("再读chan的协程结束，str=%v， ok=%v\n",str1,ok1)
    // str2,ok2 := <-cs
    // fmt.Printf("再再读chan的协程结束，str=%v， ok=%v\n",str2,ok2)

    // fmt.Println("以下是结构体chan")
    // type MyStruct struct{
    //     Name string
    // }
    // cstruct := make(chan MyStruct,3)
    // cstruct <- MyStruct{Name: "haha"}
    // close(cstruct)
    // stru,ok := <- cstruct
    // fmt.Printf("读chan的协程结束，stru=%v， ok=%v\n",stru,ok)
    // stru1,ok1 := <-cs
    // fmt.Printf("再读chan的协程结束，stru=%v， ok=%v\n",stru1,ok1)
    // stru2,ok2 := <-cs
    // fmt.Printf("再再读chan的协程结束，stru=%v， ok=%v\n",stru2,ok2)

	
    go func() {
        fmt.Print("run ")
		ticker := time.NewTicker(2 * time.Second)
    	defer ticker.Stop()

    	for range ticker.C {
        	fmt.Print("run test ")
		}

    }()
	sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
    <-sigchan

   	fmt.Print("Shutting down connectors...")

	// var a aa = aa{
	// 	aaa:&bbb{b:100,},
	// }
    // fmt.Print(a.aaa.b)


	// var p1 Person = Person{"Alice", 25}
    // fmt.Println("p1:", p1) // {Alice 25}
      
  	// var p2 *Person = &Person{"Bob", 30}
    // fmt.Println("p2:", p2)   // &{Bob 30}
    // fmt.Println("*p2:", *p2) // {Bob 30}
    // fmt.Println((*p2).Name)

} 
// type Person struct {
//     Name string
//     Age  int
// }
// 	type aa struct{
// 		aaa  *bbb 
// 	}

// 	type bbb struct{
// 		b int
// 	}