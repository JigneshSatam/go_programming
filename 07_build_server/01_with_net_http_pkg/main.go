package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	li, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Panic(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		printError(err)
		// i, err := conn.Read()
		// printError(err)
		io.WriteString(conn, "\n Hello from TCP server\n")
		fmt.Fprintln(conn, "How is your day?")
		fmt.Fprintf(conn, "%v", "Well, I hope!")

		conn.Close()
	}
}

func printError(err interface{}) {
	if err != nil {
		log.Println(err)
	}
}
