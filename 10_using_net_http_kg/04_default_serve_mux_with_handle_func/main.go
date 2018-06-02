package main

import (
	"io"
	"net/http"
)

func d(rw http.ResponseWriter, r *http.Request) {
	io.WriteString(rw, "Doggy doggy doggy")
}

func c(rw http.ResponseWriter, r *http.Request) {
	io.WriteString(rw, "Kitty kitty kitty")
}

func main() {
	http.HandleFunc("/dog/", d)
	http.HandleFunc("/cat", c)
	http.ListenAndServe(":8000", nil)
}
