
package main

import (
	"fmt"
	"reflect"
	"strings"
	"bytes"
    
)
/*
在 Go 语言中，字符串连接有多种方式，按性能从高到低排列如下：
*/
func main() {
//1 strings.Builder 零内存分配

var builder strings.Builder
builder.WriteString("Hello")
builder.WriteString(" ")
builder.WriteString("World")
result1 := builder.String()
fmt.Println(result1)


//bytes.Buffer  内存分配少
var buffer bytes.Buffer
buffer.WriteString("Hello")
buffer.WriteString(" ")
buffer.WriteString("World")
result2 := buffer.String()
fmt.Println(result2)


//切片 + string() 一次分配
bytes := make([]byte, 0, 20) // 预分配容量
fmt.Println(reflect.TypeOf(bytes))
bytes = append(bytes, "Hello"...)
bytes = append(bytes, " "...)
bytes = append(bytes, "World"...)
result3 := string(bytes)
fmt.Println(result3)

//fmt.Sprintf 有反射开销，性能一般
result4 := fmt.Sprintf("%s %s", "Hello", "World")
fmt.Println(result4)

//+ 操作符（少量拼接） 大量拼接时性能差
result5 := "Hello" + " " + "World"
fmt.Print(result5)

//strings.Join
result6 := strings.Join([]string{"Hello", " ", "World"}, "")
fmt.Println(result6)
}