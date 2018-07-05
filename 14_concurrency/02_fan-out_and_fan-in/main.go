package main

import (
	"fmt"
	"sync"
)

func main() {
	c := producer()
	parallelProcessesCount := 100
	arr := make([]chan int, parallelProcessesCount)
	for i := 1; i <= parallelProcessesCount; i++ {
		arr[i-1] = fanout(c)
	}
	out := fanIn(arr...)
	y := 0
	for n := range out {
		y++
		fmt.Println(y, "\t", n)
	}
}

func fanout(n chan int) chan int {
	c := make(chan int)
	go func() {
		for x := range n {
			c <- fact(x)
		}
		close(c)
	}()
	return c
}

func producer() chan int {
	c := make(chan int)
	go func() {
		// for i := 1; i <= 100000; i++ {
		for i := 1; i <= 100; i++ {
			for j := 1; j <= 65; j++ {
				c <- j
			}
		}
		close(c)
	}()
	return c
}

func fanIn(channels ...chan int) chan int {
	out := make(chan int)
	ws := sync.WaitGroup{}
	ws.Add(len(channels))
	go func() {
		for _, c := range channels {
			go func(c chan int) {
				for x := range c {
					out <- x
				}
				ws.Done()
			}(c)
		}
	}()
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
