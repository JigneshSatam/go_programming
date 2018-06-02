package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", dog)
	http.HandleFunc("/toby.jpeg", dogPic)
	http.ListenAndServe(":8000", nil)
}

func dog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `
		<img src="/toby.jpeg">
		`)
}

func dogPic(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("toby.jpeg")
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	io.Copy(w, f)
}
