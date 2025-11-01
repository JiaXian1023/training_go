package main

import (
    "bytes"
    "fmt"
    "strings"
    "time"
)
/*
詳細效能數據
方法					執行時間	記憶體分配	  
strings.Builder			~150μs		1次	 
bytes.Buffer			~200μs		1次	 
slice + strings.Join	~300μs	 	2次 
+ 操作符				  ~15ms		10000次	 
fmt.Sprintf				~150ms	   20000+次	 

對於1萬次字串拼接：

首選：strings.Builder（預分配記憶體）

次選：bytes.Buffer 或 slice + strings.Join

絕對避免：循環中使用 + 或 fmt.Sprintf

效能差異可能達到 100-1000倍，在高效能場景中選擇正確的拼接方法至關重要！
*/
const (
    count    = 10000
    testStr  = "a"
)

func main() {
    // 測試各種拼接方法
    fmt.Printf("拼接 %d 個字串效能測試:\n", count)
    
    // 1. strings.Builder
    start := time.Now()
    var builder strings.Builder
    builder.Grow(count * len(testStr)) // 預分配空間
    for i := 0; i < count; i++ {
        builder.WriteString(testStr)
    }
    _ = builder.String()
    fmt.Printf("strings.Builder: %v\n", time.Since(start))
    
    // 2. bytes.Buffer
    start = time.Now()
    var buffer bytes.Buffer
    buffer.Grow(count * len(testStr))
    for i := 0; i < count; i++ {
        buffer.WriteString(testStr)
    }
    _ = buffer.String()
    fmt.Printf("bytes.Buffer: %v\n", time.Since(start))
    
    // 3. slice + strings.Join
    start = time.Now()
    slice := make([]string, count)
    for i := 0; i < count; i++ {
        slice[i] = testStr
    }
    _ = strings.Join(slice, "")
    fmt.Printf("slice + strings.Join: %v\n", time.Since(start))
    
    // 4. + 操作符
    start = time.Now()
    result := ""
    for i := 0; i < count; i++ {
        result += testStr
    }
    _ = result
    fmt.Printf("+ 操作符: %v\n", time.Since(start))
    
    // 5. fmt.Sprintf
    start = time.Now()
    result = ""
    for i := 0; i < count; i++ {
        result = fmt.Sprintf("%s%s", result, testStr)
    }
    _ = result
    fmt.Printf("fmt.Sprintf: %v\n", time.Since(start))
}