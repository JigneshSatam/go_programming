package books

import (
	"fmt"
	"net/http"

	"github.com/JigneshSatam/go_programming/13_CRUD/config"
)

type book struct {
	Id   int
	Name string
}

func ShowBook(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT books.id, books.name FROM books")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// var allName string
	var bk = book{}
	for rows.Next() {
		rows.Scan(&bk.Id, &bk.Name)
	}
	// fmt.Println(allName)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Println(config.Template.DefinedTemplates())
	config.Template.ExecuteTemplate(w, "show.html", bk)
}
