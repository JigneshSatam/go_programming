package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
	var i int
	var requestMethod string
	var requestUri string
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if i == 0 {
			words := strings.Fields(line)
			fmt.Println("$$$$$$$$ Request Method: ", words[0], "  $$$$$$$$")
			fmt.Println("$$$$$$$$ Request Uri: ", words[1], "  $$$$$$$$")
			requestMethod = words[0]
			requestUri = words[1]
		}
		if line == "" {
			fmt.Println("This is end of HTTP Request Headers")
			break
		}
		i++
	}
	fmt.Println("Code got here.")
	body := "Request Method: <b>" + requestMethod + "</b>\r\n</br>"
	body += "Request Uri:&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<b>" + requestUri + "</b>\r\n</br>"
	body += "Check out the response headers and body in browser."
	io.WriteString(conn, "HTTP/1.1 200 Ok\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	io.WriteString(conn, "Content-type:text/html\r\n")
	io.WriteString(conn, "\r\n")
	io.WriteString(conn, "\n"+body+"\r\n")
	conn.Close()
}
