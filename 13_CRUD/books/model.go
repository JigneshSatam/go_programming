package books

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/JigneshSatam/go_programming/13_CRUD/config"
)

// Book is book structure in database
type Book struct {
	ID      int       `db:"id"`
	Name    string    `db:"name"`
	Release time.Time `db:"release_date" json:"-"`
	Authors []string  `db:"authors" json:"-"`
	Numbers []int     `db:"numbers" json:"-"`
}

// Books is an array of book structure
type Books []Books

// var Books []Book

func find(id int) (Book, bool) {
	rows, err := config.DB.Query("SELECT * FROM books where id = $1 limit 1", id)
	config.ParseError(err)
	book := Book{}
	for rows.Next() {
		// auth := pq.Array(&book.Authors)
		// time := pq.NullTime{Time: book.Release, Valid: false}
		// err := rows.Scan(&book.ID, &book.Name, &time, auth)
		err := config.ScanToStruct(rows, &book)
		fmt.Println("Book ==> ", book)
		config.ParseError(err)
	}
	return book, true
}

func findByName(name string) (Book, bool) {
	rows, err := config.DB.Query("SELECT * FROM books where name = $1 limit 1", name)
	config.ParseError(err)
	book := Book{}
	for rows.Next() {
		err := config.ScanToStruct(rows, &book)
		config.ParseError(err)
	}
	return book, true
}

func findAllNew() ([]Book, bool) {
	rows, err := config.DB.Query("SELECT * FROM books order by id")
	config.ParseError(err)
	bks := make([]Book, 0)
	for rows.Next() {
		bk := Book{}
		config.ScanToStruct(rows, &bk)
		bks = append(bks, bk)
	}
	return bks, true
}

func findAll() ([]Book, bool) {
	// t := reflect.TypeOf(Book{})
	// t.FieldByName("Name").Tag.Get("db")
	rows, err := config.DB.Query("SELECT * FROM books order by id")
	config.ParseError(err)
	// cols, _ := rows.Columns()
	// colTypes, _ := rows.ColumnTypes()
	// fmt.Println("cols ", cols)
	// fmt.Println("colTypes ", &colTypes[0])
	// books := make([]Book, 0)
	// for rows.Next() {
	// 	bk := Book{}
	// 	rows.Scan(&bk)
	// 	books = append(books, bk)
	// }
	// return books, true
	cols, err := rows.Columns() // Remember to check err afterwards
	// for i, _ := range cols {
	// 	vals[i] = new(sql.RawBytes)
	// }

	// fmt.Println(bk.ID)
	colIndexMapper := make(map[string]int)
	for i, col := range cols {
		colIndexMapper[col] = i
	}
	// fmt.Println(colIndexMapper)
	books := make([]Book, 0)
	vals := make([]interface{}, len(cols), len(cols))
	for rows.Next() {
		// valsptr := make([]interface{}, len(cols))
		for i, _ := range cols {
			vals[i] = &vals[i]
		}
		err = rows.Scan(vals...)

		// val := vals[0].(interface{})
		// for i, val := range vals {
		// 	x := val.(interface{})
		// 	fmt.Printf("%v %v\n", cols[i], x)
		// }
		bk := Book{}
		rv := reflect.ValueOf(&bk).Elem()
		for i := 0; i < rv.NumField(); i++ {
			f := rv.Field(i)

			tag := rv.Type().Field(i).Tag.Get("db")
			index := colIndexMapper[tag]
			value := vals[index]
			fmt.Println("New bk ", value.(interface{}))

			if f.CanSet() {
				switch value.(interface{}).(type) {
				case int64:
					x := value.(int64)
					f.SetInt(x)
				case string:
					x := value.(string)
					f.SetString(x)
				}
			}
			// tg := f.Tag.Get("db")
			// val := rv
			// fmt.Printf("FieldTag: %v,  Value: %v", tg, val)
		}
		// fmt.Println(books)
		fmt.Println("Old bk ", bk)
		books = append(books, bk)

		// Now you can check each element of vals for nil-ness,
		// and you can use type introspection and type assertions
		// to fetch the column into a typed variable.
	}
	return books, true
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
	res, err := config.DB.Exec("UPDATE books SET "+setString+" where id = $1", book.ID)
	config.ParseError(err)
	return res
}
