package main

import (
	"net/http"

	"github.com/JigneshSatam/go_programming/13_CRUD/books"
)

func main() {
	http.HandleFunc("/book", books.ShowBook)
	http.ListenAndServe(":3000", nil)

}
