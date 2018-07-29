package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	stTime := time.Now()
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
	fmt.Println("Final time taken ==> ", time.Since(stTime))
}

func makeCall() {
	stTime := time.Now()
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:3000/books/1", nil)
	if err != nil {
		log.Panicln(err)
	}
	req.SetBasicAuth("ANALYST", "ANALYST")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	fmt.Println("ct ==>", req.Header.Get("Content-Type"))
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var resJson map[string]interface{}
	x := &resJson
	err = json.Unmarshal(body, x)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Printf("Response : %v ==> %T\n", resJson["message"], resJson["message"])
	fmt.Println("Time taken", time.Since(stTime))
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
