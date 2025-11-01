package main

import (
	"fmt"
	// "strings"
	"time"
)

func app() func(string) string {
	t := "Hi"
	c := func(b string) string {
		t = t + "" + b
		return t
	}
	return c
}

func main() {
	//閉包1
	// a:=app()
	// b:=app()
	// fmt.Println(a("go"))
	// fmt.Println(b("ALL"))

	a := 5
	func() {
	    fmt.Println("a =", a)
	}()
	a = 10
	time.Sleep(1 * time.Second)

	done := make(chan bool)

	values := []string{"a", "b", "c"}
	for _, v := range values {
	    go func() {
	        fmt.Println(v)
	        done <- true
	    }()
	}
	// wait for all goroutines to complete before exiting
	for _ = range values {
		fmt.Println("?")
	    <-done
	}

	// 通过以上的讲解，对闭包应该有了更清晰的认识。如果面试中再被问到闭包，你可以这么回答：
	// 对闭包来说，函数在该语言中得是一等公民。一般来说，一个函数返回另外一个函数，这个被返回的函数可以引用外层函数的局部变量，
	// 这形成了一个闭包。通常，闭包通过一个结构体来实现，它存储一个函数和一个关联的上下文环境。但 Go 语言中，匿名函数就是一个闭包，
	// 它可以直接引用外部函数的局部变量，因为 Go 规范和 FAQ 都这么说了。

	//是一种引用了外部变量的函数
	//闭包就是一种阻止垃圾回收器将变量从内存中移除的方法，使创建变量的执行环境外面可以访问到该创建的变量

	// c1 := counter()
	// fmt.Println(c1()) // 1
	// fmt.Println(c1()) // 2

	// c2 := counter()
	// fmt.Println(c2()) // 1
	// fmt.Println(c2()) // 2
	// fmt.Println(c2()) // 3

	//  fmt.Println(c1()) // 3

	/*
		var conditions []string
		//fmt.Printf("Datatype of i : %T\n", conditions)

		fmt.Println(reflect.TypeOf(conditions).String())
		//fmt.Println(conditions.type)
		var a = []int {1,2,3,4,5}
		b := a[0:]//index = 2 的位置, 4 表示為結尾在 index = 4(不包含)
		//切的時候cap計算是從原本最大-x  (x:y)
		var c = []int {1,4,5,6,7}
		fmt.Println(len(b),cap(b),b)
		//a=c
		slicePtr := &a
		*slicePtr =c
		fmt.Println(slicePtr)
		//fmt.Println(c,cap(c))
		fmt.Println(a)
		fmt.Println(b)*/

   	// addQuery := []string{}
	// addQuery = append(addQuery, fmt.Sprintf(`"account_status" = %v`, "1"))
	// addQuery = append(addQuery, fmt.Sprintf(`"account_status_lock_time" = %v`, "2"))
	// fmt.Println(addQuery[0])
	// query := strings.Join(addQuery, " and ")
	// fmt.Println(query)

	// var arrLazy = [...]int{5, 6, 7, 8, 22}
	// fmt.Println(arrLazy)
	// fmt.Printf("Datatype of i : %T\n", arrLazy)
	// var arrLazy2 = [5]int{5, 6, 7, 8, 22}
	// fmt.Println(arrLazy2)
	// fmt.Printf("Datatype of i : %T\n", arrLazy2)

	// var urls = []string{
	// 	"http://www.google.com/",
	// 	"http://golang.org/",
	// 	"http://blog.golang.org/",
	// }
	// fmt.Printf("Datatype of i : %T\n", urls)
	// a := [...]string{"a", "b", "c", "d"}
	// for i := range a {
	// 	fmt.Println("Array item", i, "is", a[i])
	// }


	// var arrLazy = []int{5, 6, 7, 8, 22}
	// fmt.Println(arrLazy)
	// fmt.Printf("Datatype of i : %T\n", arrLazy)
	 
}

// 闭包如何引用外部变量
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}
