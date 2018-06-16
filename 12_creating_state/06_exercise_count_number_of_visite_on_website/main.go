package main

import (
	"io"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", incrCount)
	http.ListenAndServe(":8000", nil)
}

func incrCount(w http.ResponseWriter, r *http.Request) {
	var count int
	oldCookie, err := r.Cookie("Count")
	if err == http.ErrNoCookie {
		count = 1
	} else {
		count, _ = strconv.Atoi(oldCookie.Value)
		count++
	}
	countStr := strconv.Itoa(count)
	cookie := http.Cookie{
		Name:  "Count",
		Value: countStr,
	}
	http.SetCookie(w, &cookie)
	io.WriteString(w, "Number of visits: "+countStr)
}
