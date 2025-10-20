package main

import(
	"fmt"
)

func init(){
  fmt.Println("1111")

}

func init(){
  fmt.Println("2222")

}

type Person struct{
	name string
}


type Myint int

func main(){
 
 var i int=1
 var j Myint = Myint(i)
 fmt.Println(j)
 p:=&Person{name:"test"}
fmt.Println(p.name)
//fmt.Println((*p).name) //取指針裡的值 ok
//fmt.Println((&p).name) //對指針在取地址 錯誤
}