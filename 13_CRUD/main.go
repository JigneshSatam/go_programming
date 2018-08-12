package main

import (
	"log"
	_ "net/http/pprof"

	"github.com/JigneshSatam/go_programming/13_CRUD/routes"
)

func main() {
	// http.Handle("/", http.RedirectHandler("/books", http.StatusSeeOther))
	// http.HandleFunc("/books", books.Index)
	// http.HandleFunc("/books/create", books.Create)
	// http.HandleFunc("/books/new/", books.New)
	// http.HandleFunc("/books/", books.Show)
	// http.HandleFunc("/books/edit/", books.Edit)
	// http.HandleFunc("/books/update/", books.Update)
	// http.HandleFunc(`/books/delete/`, books.Delete)
	// http.HandleFunc("/favicon.ico", http.NotFound)
	// http.ListenAndServe(":3000", nil)
	log.Println(routes.StartServer())
}
