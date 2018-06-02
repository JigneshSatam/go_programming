package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listner, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalln(err)
	}
	defer listner.Close()

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println(err)
		}
		go serveRequest(conn)
	}
}

func serveRequest(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if line == "" {
			fmt.Println("This is end of HTTP Request Headers")
			break
		}
	}
	fmt.Println("Code got here.")
	body := "Check out the response headers and body in browser."
	io.WriteString(conn, "HTTP/1.1 200 Ok\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	io.WriteString(conn, "Content-type:text/plain\r\n")
	io.WriteString(conn, "\n"+body+"\r\n")
	conn.Close()
}
