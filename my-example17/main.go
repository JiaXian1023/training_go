
package main

import (
	"fmt" 
    
) 
func main() {
/*
	// 宣告並初始化陣列 A (其類型為 [3]int，長度是類型的一部分)
    arrayA := [3]int{10, 20, 30}
    
    // 1. 陣列賦值 (值複製)
    // 將 arrayA 賦值給 arrayB。由於陣列是值類型，
    // Go 會複製 arrayA 中的所有元素到 arrayB 的新記憶體空間。
    arrayB := arrayA 

    fmt.Println("--- 初始狀態 ---")
    fmt.Printf("Array A: %v (記憶體位址: %p)\n", arrayA, &arrayA)
    fmt.Printf("Array B: %v (記憶體位址: %p)\n", arrayB, &arrayB)
    fmt.Println("Array A 和 Array B 的記憶體位址不同，證明是獨立的物件。\n")

    // 2. 進行修改
    // 修改 arrayB 中的第一個元素
    arrayB[0] = 99 

    fmt.Println("--- 修改 arrayB[0] 之後 ---")
    fmt.Printf("Array A: %v\n", arrayA) // Array A 保持不變
    fmt.Printf("Array B: %v\n", arrayB) // Array B 被修改為 [99 20 30]

    fmt.Println("\n結論：修改 Array B 並不影響 Array A，因為 Array B 是 Array A 的獨立副本。")
*/


// 1. 宣告一個底層陣列 [10 20 30 40 50]
underlyingArray := [5]int{10, 20, 30, 40, 50}

// 2. 創建一個切片 s1，引用 underlyingArray 的前三個元素,切片會關聯到指標
s1 := underlyingArray[0:3] // s1: [10 20 30]

// 3. 創建另一個切片 s2，引用 underlyingArray 的中間三個元素,切片會關聯到指標
s2 := underlyingArray[1:4] // s2: [20 30 40]

// 4. 通過 s1 改變底層陣列的資料
s1[1] = 99 // 改變了底層陣列的第二個元素 (20 -> 99)

// 5. 檢查 s2 的值
// 因為 s2 也引用了底層陣列，所以它受 s1 改變的影響！
fmt.Println(s1) // 輸出: [10 99 30]
fmt.Println(s2) // 輸出: [99 30 40]  (第二個元素變成了 99)
fmt.Println(underlyingArray) // ,切片會關聯到指標所以改變了array的值

}