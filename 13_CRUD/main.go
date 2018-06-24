package main

import (
	"net/http"

	"github.com/JigneshSatam/go_programming/13_CRUD/books"
)

func main() {
	http.Handle("/", http.RedirectHandler("/books", http.StatusSeeOther))
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			books.Create(w, r)
			return
		}
		books.Index(w, r)
	})
	http.HandleFunc("/books/new/", books.New)
	http.HandleFunc("/books/", books.Show)
	http.HandleFunc(`/books/delete/`, books.Delete)
	http.ListenAndServe(":3000", nil)
}
