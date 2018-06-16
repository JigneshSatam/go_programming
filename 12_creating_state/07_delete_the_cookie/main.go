package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/set", set)
	http.HandleFunc("/read", read)
	http.HandleFunc("/expire", expire)
	http.ListenAndServe(":8000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, `
		<h3>Set Cookie</h3>
		</br>
		<button onClick="window.location='/set'">Set Cookie</button>
	`)
}

func set(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:  "Yo",
		Value: "Dude",
	}
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, `
		<h3>Your cookie is set</h3>
		</br>
		<button onClick="window.location='/read'">Read Cookie</button>
		`)
}

func read(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Yo")
	if err != nil {
		http.Redirect(w, r, "/set", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintln(w, `
		<h3>Your Cookie is: `+cookie.Value+`</h3>
		</br>
		<button onClick="window.location='/expire'">Expire Cookie</button>
		`)
}

func expire(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Yo")
	if err != nil {
		http.Redirect(w, r, "/set", http.StatusSeeOther)
		return
	}
	cookie.MaxAge = 10
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/index", http.StatusSeeOther)
}
