package main

import (
	"net/http"
	"text/template"
)

var tlp *template.Template

func init() {
	tlp = template.Must(template.New("").ParseGlob("home.gohtml"))
}

func home(rw http.ResponseWriter, r *http.Request) {
	tlp.ExecuteTemplate(rw, "home.gohtml", nil)
}

func dog(rw http.ResponseWriter, r *http.Request) {
	tlp.ExecuteTemplate(rw, "home.gohtml", nil)
}

func me(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tlp.ExecuteTemplate(rw, "home.gohtml", r.Form)
}
func main() {
	http.Handle("/", http.HandlerFunc(home))
	http.Handle("/dog", http.HandlerFunc(dog))
	http.Handle("/me/", http.HandlerFunc(me))
	http.ListenAndServe(":8000", nil)
}
