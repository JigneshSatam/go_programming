package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type user struct {
	Username string
	Password string
	Name     string
}

var sessions = map[string]string{}
var userMapping = map[string]user{}
var temp *template.Template

func init() {
	temp = template.Must(template.New("").ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", login)
	// http.Handle("/login", http.RedirectHandler("/", 301))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/create_user", createUser)
	http.HandleFunc("/socket", socket)
	http.ListenAndServe(":8000", nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
	w.Header().Set("Content-Type", "text/html")
	temp.ExecuteTemplate(w, "signup.html", nil)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		n := r.FormValue("name")
		un := r.FormValue("username")
		p := r.FormValue("password")
		userMapping[un] = user{
			Name:     n,
			Username: un,
			Password: p,
		}
		setSession(w, un)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/signup", http.StatusSeeOther)
}

func login(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	if r.Method == "POST" {
		un := r.FormValue("username")
		p := r.FormValue("password")
		if u, ok := userMapping[un]; ok {
			if u.Password == p {
				setSession(w, un)
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				return
			}
		}
		http.Error(w, "Username and/or Password is wrong", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	temp.ExecuteTemplate(w, "signin.html", nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if isLoggedIn(r) {
			cookie, err := r.Cookie("SID")
			if err == nil {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		var u user
		u, ok := getUser(r)
		if ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			temp.ExecuteTemplate(w, "dashboard.gohtml", u)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func socket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	fmt.Fprintf(w, "data: %v\n\n", time.Now())
	fmt.Printf("data: %v\n", time.Now())
}
