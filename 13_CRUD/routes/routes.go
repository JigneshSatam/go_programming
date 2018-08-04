package routes

import "github.com/JigneshSatam/go_programming/13_CRUD/books"

// Routes are multidimetional array which contails http Method string, route string and
// http HandleFunc as elements.
var Routes = [][]interface{}{
	{"GET", "/", books.Index},
	{"GET", "/books", books.Index},
	{"GET", "/books/new", books.New},
	{"GET", "/books/[0-9]+", books.Show},
	{"GET", "/books/[0-9]+/edit", books.Edit},
	{"GET", "/books/[0-9]+/delete", books.Delete},
	{"POST", "/books/[0-9]+", books.Update},
	{"POST", "/books", books.Create},
}
