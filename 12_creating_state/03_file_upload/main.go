package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", fileUpload)
	http.ListenAndServe(":8000", nil)
}

func fileUpload(w http.ResponseWriter, r *http.Request) {
	var s string
	fmt.Println("Method: ", r.Method)
	if r.Method == http.MethodPost {
		file, header, err := r.FormFile("q")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		fmt.Println("file: ", file, "\n header:", header, "\n error: ", err)
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s = string(bs)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	form := `
		<form method="post" enctype="multipart/form-data">
			Upload text file
			</br>
			</br>
			<input type="file" name="q"/>
			</br>
			</br>
			<input type="submit" />
		</form>
	`
	io.WriteString(w, form+s)
}
