package main

import (
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()
	for n := range fanIn(fanOut(producer(), 493)...) {
		fmt.Println(n)
	}
	fmt.Println("sub ==>", time.Since(startTime).Seconds())
}

func producer() chan int {
	c := make(chan int)
	go func() {
		for i := 0; i < 100000; i++ {
			for j := 0; j < 21; j++ {
				c <- j
			}
		}
		close(c)
	}()
	return c
}

func fanOut(c chan int, noOfProcesses int) []chan int {
	// arrChan := make([]chan int, noOfProcesses)
	var arrChan []chan int
	for i := 0; i < noOfProcesses; i++ {
		newC := make(chan int)
		// arrChan[i] = newC
		arrChan = append(arrChan, newC)
		go func(i int) {
			for n := range c {
				newC <- fact(n)
			}
			close(newC)
		}(i)
	}
	return arrChan
}

func fanIn(channels ...chan int) chan int {
	c := make(chan int)
	complete := make(chan bool)
	n := len(channels)
	for i := 0; i < n; i++ {
		go func(i int) {
			for x := range channels[i] {
				c <- x
			}
			complete <- true
		}(i)
	}
	go func() {
		for i := 0; i < n; i++ {
			<-complete
		}
		close(c)
	}()
	return c
}

func fact(n int) int {
	fact := 1
	for i := n; i > 1; i-- {
		fact *= i
	}
	return fact
}
