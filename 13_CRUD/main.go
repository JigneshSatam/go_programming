package main

import (
	"fmt"

	"github.com/JigneshSatam/go_programming/13_CRUD/config"
)

type book struct {
	id   int
	name string
}

func main() {
	// http.HandleFunc("")
	// http.ListenAndServe(":3000", nil)
	rows, err := config.DB.Query("SELECT books.name FROM books")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var allName string
	for rows.Next() {
		bk := book{}
		rows.Scan(&bk.name)
		allName = allName + " " + bk.name
	}
	fmt.Println(allName)
}
