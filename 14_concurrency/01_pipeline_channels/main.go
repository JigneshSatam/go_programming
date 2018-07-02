package main

import "fmt"

// Find factorial of a number
func main() {
	arr := make([]int64, 66)
	for i := 1; i <= 66; i++ {
		arr[i-1] = int64(i)
	}
	c := allFacts(arr...)
	for n := range c {
		fmt.Println(<-n)
	}
}

func allFacts(n ...int64) chan chan int64 {
	c := make(chan chan int64)
	go func(n ...int64) {
		for _, v := range n {
			c <- fact(v)
		}
		close(c)
	}(n...)
	return c
}

func fact(n int64) chan int64 {
	c := make(chan int64)
	go func(n int64) {
		fact := int64(1)
		for ; n > 0; n-- {
			fact = fact * n
		}
		c <- fact
		close(c)
	}(n)
	return c
}
