package main

import (
	"fmt"
	"sort"
)

type Pay struct {
	money int
	name  string
}

func main() {
	var a []Pay
	a = append(a, Pay{name: "john", money: 0})
	a = append(a, Pay{name: "john2", money: 322})
	a = append(a, Pay{name: "john3", money: 55})
	a = append(a, Pay{name: "john4", money: 1})
	a = append(a, Pay{name: "john5", money: 100})

	// 方法 1: 使用 sort.Slice 按 money 降序排序
	sort.Slice(a, func(i, j int) bool {
		return a[i].money > a[j].money // 降序排序
	})

	fmt.Println("按 money 降序排序:")
	for _, p := range a {
		fmt.Printf("name: %s, money: %d\n", p.name, p.money)
	}
}
