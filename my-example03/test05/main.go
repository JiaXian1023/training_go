package main

import (
	"fmt"
	"sort"
)

type Pay struct {
	name string
	top  int
	sort int
}

func main() {
	var a []Pay
	a = append(a, Pay{name: "john", top: 0, sort: 0})
	a = append(a, Pay{name: "john2", top: 4, sort: 2})
	a = append(a, Pay{name: "john3", top: 1, sort: 2})
	a = append(a, Pay{name: "john4", top: 4, sort: 1})
	a = append(a, Pay{name: "john5", top: 1, sort: 1})

	// 一行搞定：先按 top 降序，再按 sort 降序
	sort.Slice(a, func(i, j int) bool {
		//fmt.Println("i", i)
		//fmt.Println("j", j)
		//top不等於,就可以排序大的在前面
		if a[i].top != a[j].top {
			fmt.Println("a[i].top != a[j].top", i, j)
			return a[i].top > a[j].top
		}
		//top等於就看sort 大的在前面
		return a[i].sort > a[j].sort
	})

	// 輸出結果
	for _, p := range a {
		fmt.Printf("name: %s, top: %d, sort: %d\n", p.name, p.top, p.sort)
	}
}
