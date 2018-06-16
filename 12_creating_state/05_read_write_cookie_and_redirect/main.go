package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", redirect)
	http.HandleFunc("/writeCookie", setCookie)
	http.HandleFunc("/readCookie", getCookie)
	http.ListenAndServe(":8000", nil)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/writeCookie", http.StatusPermanentRedirect)
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:  "Yo",
		Value: "Dude",
	}
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `
		<h1>Cookie Set ckeck browser</h1>
		</br>
		<form method="GET" action="/readCookie">
			<input type="submit" value="See Cookie" />
		</form>
	`)
}

func getCookie(w http.ResponseWriter, r *http.Request) {
	cookiePointer, err := r.Cookie("Yo")
	if err != nil {
		log.Panic(err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w,
		`<h1>Cookie is: `+
			cookiePointer.Value+
			`</h1>
		</br>
		<input type="submit" value="Set Cookie" onClick="window.location='/writeCookie'"/>
		`)
}
