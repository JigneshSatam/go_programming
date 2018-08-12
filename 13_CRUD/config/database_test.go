package config

import (
	"log"
	"testing"
	"time"
	// It is used to have pg driver
	_ "github.com/lib/pq"
)

type Book struct {
	ID      int       `db:"id"`
	Name    string    `db:"name"`
	Release time.Time `db:"release_date" json:"release_date"`
	Authors []string  `db:"authors" json:"authors"`
	Numbers []int     `db:"numbers" json:"-"`
}

func ScanToStruct10(b *testing.B) {
	rows, err := DB.Query("SELECT * FROM books where name ilike $1 limit 1", "Jungle%")
	ParseError(err)
	book := Book{}
	x := rows.Next()
	log.Println(x)
	if x {
		for n := 0; n < b.N; n++ {
			ScanToStruct(rows, book)
		}
	}

}

func BenchmarkScanToStruct(b *testing.B) {
	rows, err := DB.Query("SELECT * FROM books where name ilike $1 limit 1", "%Jungle%")
	ParseError(err)
	book := Book{}
	if rows.Next() {
		for n := 0; n < b.N; n++ {
			ScanToStruct(rows, &book)
		}
	}

}

func BenchmarkScanToStructWithRefl(b *testing.B) {
	rows, err := DB.Query("SELECT * FROM books where name ilike $1 limit 1", "%Jungle%")
	ParseError(err)
	book := Book{}
	if rows.Next() {
		for n := 0; n < b.N; n++ {
			ScanToStructWithRefl(rows, &book)
		}
	}

}

func BenchmarkScanToStructPQLib(b *testing.B) {
	rows, err := DB.Query("SELECT * FROM books where name ilike $1 limit 1", "%Jungle%")
	ParseError(err)
	book := Book{}
	if rows.Next() {
		for n := 0; n < b.N; n++ {
			ScanToStructPQLib(rows, &book)
		}
	}

}
