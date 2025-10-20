package main

import "fmt"

func checkValue(s int) {
	switch s {
	case 3:
		fallthrough
	case 2:
		fallthrough
	case 0, 1:
		fmt.Println("check value is", s)
	}
}

func main() {
	checkValue(3) 
	checkValue(2)
}
