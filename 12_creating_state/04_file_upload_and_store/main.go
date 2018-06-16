package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", fileUploadAndStore)
	http.ListenAndServe(":8000", nil)
}

func fileUploadAndStore(w http.ResponseWriter, r *http.Request) {
	var s string
	if r.Method == http.MethodPost {
		file, header, err := r.FormFile("q")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dstFilePointer, err := os.Create(filepath.Join("./users", header.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dstFilePointer.Close()
		_, err = dstFilePointer.Write(bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s = string(bs)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	form := `
	<form method="post" enctype="multipart/form-data">
		Upload text file</br></br>
		<input type="file" name="q"/></br></br>
		<input type="submit"/></br></br>
	</form>
	`
	io.WriteString(w, form+s)
}
