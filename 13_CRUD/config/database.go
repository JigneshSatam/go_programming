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
	migsDir := "migrations"
	mig, err := ioutil.ReadDir(migsDir)
	if len(mig) == 0 {
		migsDir = "../migrations"
		mig, err = ioutil.ReadDir(migsDir)
	}
	ParseError(err)
	fmt.Println("\n==============================================================")
	fmt.Println("                        Running Migrations")
	for i, m := range mig {
		fmt.Printf("                     %d. %s\n", i+1, m.Name())
		dat, err := ioutil.ReadFile(migsDir + "/" + m.Name())
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
			// case reflect.TypeOf(*new([]int64)), reflect.TypeOf(*new([]float64)), reflect.TypeOf(*new([]bool)), reflect.TypeOf(*new([]byte)):
			// 	pq.Array(f.Addr().Interface()).Scan(value)
			case reflect.TypeOf(*new([]string)):
				f.Set(reflect.ValueOf(getStringArr(pqArrayToArray(value.([]byte)))).Convert(fType.Type))
			case reflect.TypeOf(*new([]int64)):
				f.Set(reflect.ValueOf(getInt64Arr(pqArrayToArray(value.([]byte)))).Convert(fType.Type))
			case reflect.TypeOf(*new([]int32)):
				f.Set(reflect.ValueOf(getInt32Arr(pqArrayToArray(value.([]byte)))).Convert(fType.Type))
			case reflect.TypeOf(*new([]int16)):
				f.Set(reflect.ValueOf(getInt16Arr(pqArrayToArray(value.([]byte)))).Convert(fType.Type))
			case reflect.TypeOf(*new([]int8)):
				f.Set(reflect.ValueOf(getInt8Arr(pqArrayToArray(value.([]byte)))).Convert(fType.Type))
			case reflect.TypeOf(*new([]int)):
				f.Set(reflect.ValueOf(getIntArr(pqArrayToArray(value.([]byte)))).Convert(fType.Type))
			case reflect.TypeOf(*new([]float64)), reflect.TypeOf(*new([]bool)), reflect.TypeOf(*new([]byte)):
				arrByteArr := pqArrayToArray(value.([]byte))
				intArr := make([]int, len(arrByteArr), len(arrByteArr))
				// arr := reflect.MakeSlice(fType.Type, 0, 0)
				for i, byteArr := range arrByteArr {
					// for _, byteArr := range arrByteArr {
					if intVal, err := strconv.Atoi(string(byteArr)); err == nil {
						intArr[i] = intVal
						// arr = reflect.Append(arr, reflect.ValueOf(intVal))
					} else {
						ParseError(err)
					}
				}
				// f.Set(arr)
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

	return nil
}

// ScanToStructWithRefl is abd
func ScanToStructWithRefl(row *sql.Rows, strt interface{}) error {
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
			t := fType.Type
			arrByteArr := pqArrayToArray(value.([]byte))
			arr := reflect.MakeSlice(t, 0, 0)
			for _, byteArr := range arrByteArr {
				val, err := getValueOfBytes(byteArr, t)
				ParseError(err)
				arr = reflect.Append(arr, val)
			}
			f.Set(arr)
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

	return nil
}

// ScanToStructPQLib bacgd ede
func ScanToStructPQLib(row *sql.Rows, strt interface{}) error {
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
			// case reflect.TypeOf(*new([]int)):
			// 	arrByteArr := pqArrayToArray(value.([]byte))
			// 	// intArr := make([]int, len(arrByteArr), len(arrByteArr))
			// 	arr := reflect.MakeSlice(fType.Type, 0, 0)
			// 	for _, byteArr := range arrByteArr {
			// 		if intVal, err := strconv.Atoi(string(byteArr)); err == nil {
			// 			// intArr[i] = intVal
			// 			arr = reflect.Append(arr, reflect.ValueOf(intVal))
			// 		} else {
			// 			ParseError(err)
			// 		}
			// 	}
			// 	f.Set(arr)
			// 	// f.Set(reflect.ValueOf(intArr).Convert(fType.Type))
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

func getValueOfBytes(b []byte, t reflect.Type) (reflect.Value, error) {
	var v interface{}
	var err error
	switch t {
	case reflect.TypeOf(*new([]string)):
		v = string(b)
	case reflect.TypeOf(*new([]int)):
		v, err = strconv.Atoi(string(b))
		if err != nil {
			return reflect.ValueOf(nil), err
		}
	case reflect.TypeOf(*new([]int64)):
		v, err = strconv.Atoi(string(b))
		v = int64(v.(int))
		if err != nil {
			return reflect.ValueOf(nil), err
		}
	case reflect.TypeOf(*new([]float32)):
		v, err = strconv.ParseFloat(string(b), 32)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
	}
	return reflect.ValueOf(v), err
}

func getStringArr(arrByteArr [][]byte) []string {
	strArr := make([]string, len(arrByteArr), len(arrByteArr))
	for i, byteArr := range arrByteArr {
		strArr[i] = string(byteArr)
	}
	return strArr
}

func getIntArr(arrByteArr [][]byte) []int {
	intArr := make([]int, len(arrByteArr), len(arrByteArr))
	for i, byteArr := range arrByteArr {
		if num, err := strconv.Atoi(string(byteArr)); err == nil {
			intArr[i] = num
		}
	}
	return intArr
}

func getInt8Arr(arrByteArr [][]byte) []int8 {
	intArr := make([]int8, len(arrByteArr), len(arrByteArr))
	for i, byteArr := range arrByteArr {
		if num, err := strconv.Atoi(string(byteArr)); err == nil {
			intArr[i] = int8(num)
		}
	}
	return intArr
}

func getInt16Arr(arrByteArr [][]byte) []int16 {
	intArr := make([]int16, len(arrByteArr), len(arrByteArr))
	for i, byteArr := range arrByteArr {
		if num, err := strconv.Atoi(string(byteArr)); err == nil {
			intArr[i] = int16(num)
		}
	}
	return intArr
}

func getInt32Arr(arrByteArr [][]byte) []int32 {
	intArr := make([]int32, len(arrByteArr), len(arrByteArr))
	for i, byteArr := range arrByteArr {
		if num, err := strconv.Atoi(string(byteArr)); err == nil {
			intArr[i] = int32(num)
		}
	}
	return intArr
}

func getInt64Arr(arrByteArr [][]byte) []int64 {
	intArr := make([]int64, len(arrByteArr), len(arrByteArr))
	for i, byteArr := range arrByteArr {
		if num, err := strconv.Atoi(string(byteArr)); err == nil {
			intArr[i] = int64(num)
		}
	}
	return intArr
}
