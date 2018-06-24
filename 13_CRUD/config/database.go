package config

import (
	"database/sql"
	"fmt"
	"strings"

	// It is used to have pg driver
	_ "github.com/lib/pq"
)

// DB is used to connect go_database
var DB *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "baldor"
	password = "baldor123"
	dbname   = "postgres"
)

func init() {
	var err error
	dbAddr := `host=localhost port=5432 user=baldor password=baldor123 dbname=postgres sslmode=disable`
	DB, err = sql.Open("postgres", dbAddr)
	ParseError(err)
	res, err := DB.Exec("CREATE DATABASE " + "go_database")
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
	dbAddr2 := `host=localhost port=5432 user=baldor password=baldor123 dbname=go_database sslmode=disable`
	DB, err = sql.Open("postgres", dbAddr2)
	ParseError(err)
	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("\n==============================================================")
	fmt.Println("                 You are connected to database                  ")
	fmt.Println("==============================================================")
}
