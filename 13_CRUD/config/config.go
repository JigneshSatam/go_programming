package config

import (
	"database/sql"
	"fmt"

	// It is used
	_ "github.com/lib/pq"
)

// "DB is used to connect go_database"
var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("postgres", "postgres://baldor:baldor123@localhost:5432/go_database?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("==============================")
	fmt.Println("You are connected to database.")
	fmt.Println("==============================\n")
}
