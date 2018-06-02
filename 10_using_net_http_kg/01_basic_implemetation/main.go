package main

import (
	"fmt"
	"net/http"
)

type hotdog int

func (d hotdog) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Hey Dude.")
}
func main() {
	var d hotdog
	http.ListenAndServe(":8000", d)
}
