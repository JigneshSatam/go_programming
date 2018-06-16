package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":8000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("q")
	w.Header().Set("Content-Type", "text/html charset=UTF-8")
	io.WriteString(w,
		`
		<form method=post>
			Search here
			<input type=text name=q value=`+value+`>
			<input type=submit>
		</form>
		</br>
		`+"Search string: "+value)
}
