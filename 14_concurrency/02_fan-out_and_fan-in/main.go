package main

import (
	"fmt"
	"sync"
)

func main() {
	channels := producer()
	out := solver(channels...)
	for n := range out {
		fmt.Println(n)
	}
}

func producer() []chan int {
	var arr []chan int
	// for i := 1; i <= 100000; i++ {
	for i := 1; i <= 100; i++ {
		arr = append(arr, subProducer())
	}
	// go func() {
	// }()
	return arr
}

func subProducer() chan int {
	c := make(chan int)
	go func() {
		for j := 1; j <= 65; j++ {
			c <- j
		}
		close(c)
	}()
	return c
}

func solver(channels ...chan int) chan int {
	out := make(chan int)
	ws := sync.WaitGroup{}
	ws.Add(len(channels))
	for _, c := range channels {
		go func(n chan int) {
			for x := range n {
				out <- fact(x)
			}
			ws.Done()
		}(c)
	}
	go func() {
		ws.Wait()
		close(out)
	}()
	return out
}

func fact(n int) int {
	res := 1
	for i := n; i >= 1; i-- {
		res *= i
	}
	return res
}
