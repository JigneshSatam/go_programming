package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", dog)
	http.HandleFunc("/toby", dogPic)
	http.ListenAndServe(":8000", nil)
}

func dog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<img src="/toby">`)
}

func dogPic(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("toby.jpeg")
	if err != nil {
		log.Panic(err)
	}
	fi, err := f.Stat()
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
}
