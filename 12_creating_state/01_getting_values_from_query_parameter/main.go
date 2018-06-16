package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":8000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("q")
	io.WriteString(w, "You search for: "+value)
}
