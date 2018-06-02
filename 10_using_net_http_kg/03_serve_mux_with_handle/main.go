package main

import (
	"io"
	"net/http"
)

type hotdog int
type hotcat int

func (d hotdog) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	io.WriteString(rw, "Doggy doggy doggy")
}

func (c hotcat) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	io.WriteString(rw, "Kitty kitty kitty")
}

func main() {
	mux := http.NewServeMux()
	var d hotdog
	var c hotcat
	mux.Handle("/dog/", d)
	mux.Handle("/cat", c)
	http.ListenAndServe(":8000", mux)
}
