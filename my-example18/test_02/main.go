package main

import (
	"bytes"
	"fmt"
)

/*
反向字串
利用rune切片處理
如果要正確取得 index 與字元的話，可以先將字串轉成 []rune
再利用for 將第一個 最後一個元素交換
*/




func ReverseString(s string) string {
    // 將字串轉換為 rune slice 以正確處理 Unicode
    runes := []rune(s)
    
    // 使用雙指針進行反向
	//Hello World
	//i=0, j=10
	//0<10,i+1,j-1...>1,9,2,8...>3,7...>4,6...>5,5
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		//利用交換反轉
		fmt.Println("i,j",i,j)
        runes[i], runes[j] = runes[j], runes[i]
    }
    
    return string(runes)
}


func ReverseStringWithBuffer(s string) string {
    runes := []rune(s)
    var result bytes.Buffer
    result.Grow(len(runes)) // 預分配空間以提高效能
   	fmt.Println(len(runes))//11

	//0-10
	//10排到0
	for i:= len(runes)-1 ;i>=0;i--{
			  result.WriteRune(runes[i])
	}
    return result.String()
}

func ReverseStringRecursive(s string) string {
	//Hello World
	fmt.Println("s",s)
    runes := []rune(s)
	//fmt.Println("runes",runes)
    if len(runes) <= 1 {
        return s
    }
	fmt.Println(string(runes[0]))
	//遞迴把左邊排到右邊最後一個
    return ReverseStringRecursive(string(runes[1:])) + string(runes[0])
}

func main() {
    testCases := []string{
        "Hello World",
        // "你好，世界！",
        // "🚀🐹🌟", // Emoji 測試
        // "a",
        // "",
    }
    
    for _, test := range testCases {
        reversed := ReverseStringRecursive(test)
        fmt.Printf("原始: %q\n反向: %q\n\n", test, reversed)
    }
}