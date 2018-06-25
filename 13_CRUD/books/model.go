package books

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/JigneshSatam/go_programming/13_CRUD/config"
)

// Book is book structure in database
type Book struct {
	ID   int
	Name string
}

// Books is an array of book structure
var Books []Book

func find(id int) (Book, bool) {
	row := config.DB.QueryRow("SELECT * FROM books where id = $1", id)
	book := Book{}
	err := row.Scan(&book.ID, &book.Name)
	if err == nil {
		return book, true
	}
	return book, false
}

func findByName(name string) (Book, bool) {
	row := config.DB.QueryRow("SELECT * FROM books where name = $1", name)
	book := Book{}
	err := row.Scan(&book.ID, &book.Name)
	if err == nil {
		return book, true
	}
	return book, false
}

func all() ([]Book, bool) {
	books := make([]Book, 0)
	rows, err := config.DB.Query("SELECT * FROM books")
	cols, _ := rows.Columns()
	colNames := make([]interface{}, len(cols))
	colNamePtrs := make([]interface{}, len(cols))
	for i := 0; i < len(colNames); i++ {
		colNamePtrs[i] = &colNames[i]
	}
	defer rows.Close()
	present := false
	if err == nil {
		for rows.Next() {
			bk := Book{}
			err := rows.Scan(colNamePtrs...)
			if err != nil {
				return []Book{}, false
			}

			mappings := map[string]interface{}{}
			for i, colName := range cols {
				value := colNamePtrs[i].(*interface{})
				mappings[colName] = value
			}
			mapStr, err := json.Marshal(mappings)
			config.ParseError(err)
			err = json.Unmarshal(mapStr, &bk)
			config.ParseError(err)
			books = append(books, bk)
		}
		present = true
	}
	return books, present
}

func create(book Book) sql.Result {
	res, err := config.DB.Exec("INSERT INTO books (name) VALUES ($1)", book.Name)
	if err != nil {
		config.ParseError(err)
	}
	return res
}

func delete(book Book) sql.Result {
	res, err := config.DB.Exec("DELETE FROM books where id = $1", book.ID)
	config.ParseError(err)
	return res
}

func update(book Book) sql.Result {
	var mappingHash map[string]interface{}
	bookJSON, err := json.Marshal(book)
	config.ParseError(err)
	err = json.Unmarshal(bookJSON, &mappingHash)
	config.ParseError(err)
	setString := ""
	for key, value := range mappingHash {
		setString += fmt.Sprintf("%v='%v', ", key, value)
	}
	setString = strings.TrimSuffix(setString, ", ")
	// setString = strings.ToLower(setString)
	fmt.Println("setString   ", setString)
	res, err := config.DB.Exec("UPDATE books SET "+setString+" where id = $1", book.ID)
	config.ParseError(err)
	return res
}
