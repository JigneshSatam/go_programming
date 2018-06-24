package books

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/JigneshSatam/go_programming/13_CRUD/config"
)

// Show is used to fetch single book
func Show(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.RequestURI, "/books/")
	id, err := strconv.Atoi(idStr)
	if err == nil {
		book, ok := find(id)
		if ok {
			config.Templates["templates/books"].ExecuteTemplate(w, "show.html", book)
			return
		}
	}
	http.NotFound(w, r)
	// fmt.Println(config.Templates["templates/books"].DefinedTemplates())
	// var tmpl *template.Template
	// tmpl := template.Must(template.New("").ParseFiles("templates/books/show.html"))
	// tmpl.ExecuteTemplate(w, "show.html", FindBook(1))
}

// Index is to show all books
func Index(w http.ResponseWriter, r *http.Request) {
	books, ok := all()
	if ok {
		config.Templates["templates/books"].ExecuteTemplate(w, "index.html", books)
		return
	}
	http.NotFound(w, r)
}

func Create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	book := Book{
		Name: name,
	}
	create(book)
	book, ok := findByName(name)
	if ok {
		config.Templates["templates/books"].ExecuteTemplate(w, "show.html", book)
		return
	}
	http.Redirect(w, r, "/books", http.StatusSeeOther)
}

func New(w http.ResponseWriter, r *http.Request) {
	config.Templates["templates/books"].ExecuteTemplate(w, "new.html", nil)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.RequestURI, "/books/delete/")
	id, err := strconv.Atoi(idStr)
	if err == nil {
		book, ok := find(id)
		if ok {
			delete(book)
			http.Redirect(w, r, "/books", http.StatusSeeOther)
			return
		}
	}
	http.NotFound(w, r)
}
