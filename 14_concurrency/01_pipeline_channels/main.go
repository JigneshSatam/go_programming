package main

import "fmt"

// Find factorial of a number
func main() {

	adderChannel := addNumbers()

	printerChannel := allFacts(adderChannel)

	for n := range printerChannel {
		fmt.Println(n)
	}
}

func addNumbers() chan int64 {
	c := make(chan int64)
	go func() {
		for j := 1; j <= 200; j++ {
			for i := 1; i <= 50; i++ {
				c <- int64(i)
			}
		}
		close(c)
	}()
	return c
}

func allFacts(c chan int64) chan int64 {
	nc := make(chan int64)
	go func() {
		for n := range c {
			nc <- fact(n)
		}
		close(nc)
	}()

	return nc
}

func fact(n int64) int64 {
	fact := int64(1)
	for ; n > 0; n-- {
		fact = fact * n
	}
	return fact
}
