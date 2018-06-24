package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	// It is used to have pg driver
	_ "github.com/lib/pq"
)

// DB is used to connect go_database
var DB *sql.DB

const (
	host      = "localhost"
	port      = 5432
	user      = "baldor"
	password  = "baldor123"
	dbname    = "postgres"
	newdbname = "go_database"
)

func init() {
	var err error
	dbAddr := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", host, port, user, password, dbname)
	DB, err = sql.Open("postgres", dbAddr)
	ParseError(err)
	defer DB.Close()
	res, err := DB.Exec("CREATE DATABASE " + newdbname)
	if res != nil {
		fmt.Println(res)
	}
	if err != nil {
		if strings.Contains(err.Error(), ("database \"" + "go_database" + "\" already exists")) {
			fmt.Println("\n==============================================================")
			fmt.Println("                    Database already exists                   ")
			fmt.Println("==============================================================")
		} else {
			ParseError(err)
		}
	}
	DB.Close()
	dbAddr2 := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", host, port, user, password, newdbname)
	DB, err = sql.Open("postgres", dbAddr2)
	ParseError(err)
	if err = DB.Ping(); err != nil {
		panic(err)
	}

	dat, err := ioutil.ReadFile("migrations/01_create_table_books")
	ParseError(err)
	res, err = DB.Exec(string(dat))

	if err != nil && !strings.Contains(err.Error(), `relation "books" already exists`) {
		ParseError(err)
	}
	// fmt.Print(string(dat))

	fmt.Println("\n==============================================================")
	fmt.Println("                 You are connected to database                  ")
	fmt.Println("==============================================================")
}
