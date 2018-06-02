package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	li, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Panic(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
		}
		go serveRequest(conn)
	}
}

func serveRequest(conn net.Conn) {
	readRequest(conn)
	writeToRequest(conn)
}

func readRequest(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if i == 0 {
			words := strings.Fields(line)
			fmt.Println("^^^^Method :" + words[0] + "^^^^")
			fmt.Println("^^^^URI :" + words[1] + "^^^^")
		}
		if line == "" {
			break
		}
		i++
	}
}

func writeToRequest(conn net.Conn) {
	body := `<!DOCTYPE html>
		<html lang='en'>
			<head>
				<meta charset="UTF-8">
				<title>First Server Code</title>
			</head>
			<body>
				<h1>
					Hello World
				</h1>
			</body>
		</html>`
	fmt.Fprint(conn, "HTTP/1.1 200 Ok\r\n")
	fmt.Fprintf(conn, "Content-length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-type: text/html\r\n")
	fmt.Fprintln(conn, "")
	fmt.Fprint(conn, body)
}
