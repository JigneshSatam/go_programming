package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

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
			fmt.Println("                    Database already exists")
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

	mig, err := ioutil.ReadDir("migrations")
	ParseError(err)
	fmt.Println("\n==============================================================")
	fmt.Println("                        Running Migrations")
	for i, m := range mig {

		fmt.Printf("                     %d. %s\n", i+1, m.Name())
		dat, err := ioutil.ReadFile("migrations/" + m.Name())
		ParseError(err)
		res, err = DB.Exec(string(dat))
	}
	fmt.Println("==============================================================")

	if err != nil && !strings.Contains(err.Error(), `relation "books" already exists`) {
		ParseError(err)
	}
	// fmt.Print(string(dat))

	fmt.Println("\n==============================================================")
	fmt.Println("                 You are connected to database")
	fmt.Println("==============================================================")
}

// ScanToStruct updateds the struct with values from db scan method
func ScanToStruct(row *sql.Rows, strt interface{}) error {
	// fmt.Println("Started ==> ")
	// sTime := time.Now()
	cols, err := row.Columns()
	ParseError(err)
	colIndexMapper := make(map[string]int)
	vals := make([]interface{}, len(cols), len(cols))
	for i, col := range cols {
		colIndexMapper[col] = i
		vals[i] = &vals[i]
	}
	rv := reflect.ValueOf(strt).Elem()
	err = row.Scan(vals...)
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		fType := rv.Type().Field(i)
		tag := fType.Tag.Get("db")
		value := vals[colIndexMapper[tag]]
		if value == nil {
			continue
		}
		// fmt.Println("Value ==> ", value)
		switch f.Kind() {
		case reflect.Slice:
			switch fType.Type {
			case reflect.TypeOf(*new([]string)), reflect.TypeOf(*new([]int64)), reflect.TypeOf(*new([]float64)), reflect.TypeOf(*new([]bool)), reflect.TypeOf(*new([]byte)):
				pq.Array(f.Addr().Interface()).Scan(value)
			case reflect.TypeOf(*new([]int)):
				arrByteArr := pqArrayToArray(value.([]byte))
				intArr := make([]int, len(arrByteArr), len(arrByteArr))
				for i, byteArr := range arrByteArr {
					if intVal, err := strconv.Atoi(string(byteArr)); err == nil {
						intArr[i] = intVal
					} else {
						ParseError(err)
					}
				}
				f.Set(reflect.ValueOf(intArr).Convert(fType.Type))
			case reflect.TypeOf(*new([]float32)):
				arrByteArr := pqArrayToArray(value.([]byte))
				intArr := make([]float32, len(arrByteArr), len(arrByteArr))
				for i, byteArr := range arrByteArr {
					if intVal, err := strconv.ParseFloat(string(byteArr), 32); err == nil {
						intArr[i] = float32(intVal)
					} else {
						ParseError(err)
					}
				}
				f.Set(reflect.ValueOf(intArr).Convert(fType.Type))
			default:

			}
		case reflect.Struct:
			switch fType.Type {
			case reflect.TypeOf(time.Time{}):
				nt := pq.NullTime{}
				nt.Scan(value)
				if nt.Valid {
					f.Set(reflect.ValueOf(value).Convert(fType.Type))
				}
			}
		default:
			f.Set(reflect.ValueOf(value).Convert(fType.Type))
		}
	}
	ParseError(err)
	// fmt.Println("Time Taken ==> ", time.Since(sTime))
	// fmt.Printf("%T ==> %v\n", strt, strt)
	return nil
}

func pqArrayToArray(bytes []byte) [][]byte {
	arr := make([][]byte, 0, 0)
	if bytes[0] == '{' && bytes[len(bytes)-1] == '}' {
		stByte := 1
		edByte := 1
		for i := 1; i < len(bytes); i++ {
			if bytes[i] == ',' || bytes[i] == '}' {
				arr = append(arr, bytes[stByte:edByte])
				stByte = i + 1
				edByte = i + 1
			} else {
				edByte++
			}
		}
	}
	return arr
}
