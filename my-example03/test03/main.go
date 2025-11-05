package main

import "fmt"

// 降序氣泡排序
func bubbleSortDesc(arr []int) []int {
	n := len(arr)
	// 複製切片，避免修改原始數據
	sorted := make([]int, n)
	copy(sorted, arr)
	fmt.Println("@n", n-1)
	x := 0
	for i := 0; i < n-1; i++ {
		//執行n-1遍

		fmt.Println("i", i)
		// swapped := false
		//fmt.Println("@j", n-i-1)
		for j := 0; j < n-i-1; j++ {
			x += 1
			fmt.Println("j", j)
			//n-1遍
			// 改成 > 比較，讓大的元素往前移動
			if sorted[j] < sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
				//swapped = true
			}
		}
		// 如果這一輪沒有交換，說明已經排序完成
		// if !swapped {
		// 	break
		// }
	}
	fmt.Println("@x", x)
	return sorted
}

// 原地排序版本（會修改原始切片）
func bubbleSortDescInPlace(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			if arr[j] < arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
}

func main() {
	a := []int{1, 2, 10, 3, 5}

	// 使用複製版本
	sorted := bubbleSortDesc(a)
	fmt.Println("原始陣列:", a)       // [1 2 10 3 5]
	fmt.Println("降序排序後:", sorted) // [10 5 3 2 1]

	// 使用原地排序版本
	// b := []int{1, 2, 10, 3, 5}
	// bubbleSortDescInPlace(b)
	// fmt.Println("原地排序後:", b) // [10 5 3 2 1]
}
