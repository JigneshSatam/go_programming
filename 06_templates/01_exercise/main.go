package main

import (
	"fmt"
	"os"
	"text/template"
)

func show10() int {
	return 10
}

var myFuncs = template.FuncMap{
	"show10": show10,
}

type hotel struct {
	Name    string
	Address string
	City    string
	Zip     string
	Region  string
}

type hotelList struct {
	Hotels []hotel
}

func main() {
	fmt.Println("Hello Template")
	tpl := template.Must(template.New("").Funcs(myFuncs).ParseGlob("templates/*"))

	h1 := hotel{
		Name:    "Hotel1",
		Address: "Address1",
		City:    "City1",
		Zip:     "Zip1",
		Region:  "Northern",
	}
	h2 := hotel{
		"Hotel2",
		"Address1",
		"City1",
		"Zip1",
		"Northern",
	}

	hotels := []hotel{h1, h2}
	data := hotelList{
		hotels,
	}

	tpl.ExecuteTemplate(os.Stdout, "test.gohtml", data)
}
