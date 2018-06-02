package main

import "fmt"

func main() {
	increament1 := wrapper()
	fmt.Println(increament1())
	fmt.Println(increament1())

	increament2 := wrapper()
	fmt.Println(increament2())
	fmt.Println(increament2())
}

func wrapper() func() int {
	x := 0
	return func() int {
		x++
		return x
	}
}
