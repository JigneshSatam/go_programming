package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	listner, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Panic(err)
	}

	defer listner.Close()

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println(err)
		}

		go serverRequest(conn)
	}
}

func serverRequest(conn net.Conn) {
	methodName := readRequest(conn)
	writeToRequest(conn, methodName)
}

var routes = map[string]map[string]string{
	"GET": map[string]string{
		"/": "home",
	},
	"POST": map[string]string{
		"/post": "post",
	},
}

func readRequest(conn net.Conn) string {
	scanner := bufio.NewScanner(conn)
	i := 0
	var x string
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if i == 0 {
			words := strings.Fields(line)
			method := words[0]
			uri := words[1]
			fmt.Println("^^^^^^ method: " + method + " url: " + uri)
			x = routes[method][uri]
			fmt.Println(x)
		}
		if line == "" {
			break
		}
		i++
	}
	return x
}

func writeToRequest(conn net.Conn, methodName string) {
	var body string
	var status int
	var statusMessage string
	switch methodName {
	case "home":
		body, status, statusMessage = home()
	case "post":
		body, status, statusMessage = post()
	default:
		body, status, statusMessage = notFound()
	}

	fmt.Fprint(conn, "HTTP/1.1 "+string(status)+" "+statusMessage+" \r\n")
	fmt.Fprintf(conn, "Content-length: %d\r\n", len(body))
	fmt.Fprintln(conn)
	fmt.Fprint(conn, body)
}

func home() (string, int, string) {
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
				<form action="/post" method="post">
					<input type="submit" value="POST"/>
				</form>
			</body>
		</html>`
	return body, 200, "Ok"
}

func post() (string, int, string) {
	body := `<!DOCTYPE html>
		<html lang='en'>
			<head>
				<meta charset="UTF-8">
				<title>First Server Code</title>
			</head>
			<body>
				<h1>
					Posted successfully
				</h1>
				<a href="/">Home</a>
			</body>
		</html>`
	return body, 200, "Ok"
}

func notFound() (string, int, string) {
	body := `<!DOCTYPE html>
		<html lang='en'>
			<head>
				<meta charset="UTF-8">
				<title>First Server Code</title>
			</head>
			<body>
				<h1>
					Page Not Found
				</h1>
				<a href="/">Home</a>
			</body>
		</html>`
	return body, 404, "NotFound"
}
