package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	fmt.Println("Server is starting...")
	li, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Panic(err)
	}
	defer li.Close()

	fmt.Println("Server started and Listening on port :8000")
	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
		}
		go serveRequest(conn)
	}
}

func serveRequest(conn net.Conn) {
	err := conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Println(err)
	}
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		fmt.Fprintf(conn, "Hi from server for your: %s\n", line)
	}
	defer conn.Close()

	fmt.Println("Connection is still Open!!!")
}
