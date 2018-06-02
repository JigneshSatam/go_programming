package main

import (
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
	io.WriteString(conn, "I see you connectd\n")
	conn.Close()
}
