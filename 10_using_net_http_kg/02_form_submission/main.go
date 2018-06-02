package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").ParseGlob("*.gohtml"))
}

type hotdog int

func (d hotdog) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	tpl.ExecuteTemplate(rw, "home.gohtml", r.Form)
}
func main() {
	var d hotdog
	http.ListenAndServe(":8000", d)
}
