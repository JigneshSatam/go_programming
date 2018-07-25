package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"reflect"
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
	cols, err := row.Columns()
	ParseError(err)
	colIndexMapper := make(map[string]int)
	// tagNameMapper := make(map[string]string)
	vals := make([]interface{}, len(cols), len(cols))
	for i, col := range cols {
		colIndexMapper[col] = i
		vals[i] = &vals[i]
	}
	rv := reflect.ValueOf(strt).Elem()
	// strtIndex := make([]int, 0, len(cols))
	// for i := 0; i < rv.NumField(); i++ {
	// 	f := rv.Field(i)
	// 	structField := rv.Type().Field(i)
	// 	tag := structField.Tag.Get("db")
	// 	tagNameMapper[tag] = structField.Name
	// 	if index, ok := colIndexMapper[tag]; ok {
	// 		switch f.Type().Kind() {
	// 		case reflect.Slice:
	// 			switch f.Interface().(type) {
	// 			case []int:
	// 				newVal := f.Convert(reflect.SliceOf(reflect.TypeOf(int(0))))
	// 				f.Set(newVal)
	// 				fmt.Println("f.Type() ==> ", f.Type())
	// 				vals[index] = pq.Array(f.Addr().Interface())
	// 			default:
	// 				vals[index] = pq.Array(f.Addr().Interface())
	// 			}
	// 		case reflect.Struct:
	// 			switch structField.Type {
	// 			case reflect.TypeOf(time.Time{}):
	// 				strtIndex = append(strtIndex, index)
	// 				vals[index] = &pq.NullTime{Time: f.Interface().(time.Time), Valid: false}
	// 			}
	// 		default:
	// 			vals[index] = f.Addr().Interface()
	// 		}
	// 	} else {
	// 		fmt.Println("Tag not found ==> ", tag)
	// 	}
	// }
	err = row.Scan(vals...)
	// for _, ind := range strtIndex {
	// 	switch vals[ind].(type) {
	// 	case *pq.NullTime:
	// 		y := vals[ind].(*pq.NullTime)
	// 		if y.Valid {
	// 			val := reflect.ValueOf(y.Time)
	// 			f := rv.FieldByName(tagNameMapper[cols[ind]])
	// 			f.Set(val)
	// 		}
	// 	}
	// }

	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		fType := rv.Type().Field(i)
		tag := fType.Tag.Get("db")
		value := vals[colIndexMapper[tag]]
		switch f.Kind() {
		case reflect.Slice:
			switch fType.Type {
			case reflect.TypeOf(*new([]string)), reflect.TypeOf(*new([]int64)), reflect.TypeOf(*new([]float64)), reflect.TypeOf(*new([]bool)), reflect.TypeOf(*new([]byte)):
				pq.Array(f.Addr().Interface()).Scan(value)
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
	fmt.Printf("%T ==> %v\n", strt, strt)
	return nil

	// value := vals[index]
	// if f.CanSet() {
	// 	// val := reflect.ValueOf(value)
	// 	fieldType := f.Type()
	// 	if i == 2 {
	// 		// o := reflect.TypeOf(val)
	// 		// r := f.Type()

	// 		// x := interface{}
	// 		// x := unsafe.Pointer(2)
	// 		a := new([]string)
	// 		b := *a
	// 		// b := (*([]string))
	// 		// fmt.Printf("F ptr %v ==> %v\n", f.Type(), f.Addr().Pointer())
	// 		// b := new([]string)
	// 		// x := make([]interface{}, 1)
	// 		// x[0] =
	// 		// y := (new(interface{}))
	// 		fmt.Printf("F Type %T ==> %v\n", string(value.([]uint8)), string(value.([]uint8)))
	//	  // pq.Array(&b).Scan(value)
	// 		// x, err := fmt.Sscan(",", &b)
	// 		// ParseError(err)
	// 		// fmt.Println(x)
	// 		// val := reflect.ValueOf(*y)
	// 		val := reflect.ValueOf(b)
	// 		f.Set(val.Convert(fieldType))
	// 	} else {
	// 		val := reflect.ValueOf(value)
	// 		f.Set(val.Convert(fieldType))
	// 	}
	// }
}
