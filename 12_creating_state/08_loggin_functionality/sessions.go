package main

import (
	"net/http"
)

func isLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("SID")
	if err == http.ErrNoCookie {
		return false
	}

	if un, ok := sessions[cookie.Value]; ok {
		if _, ok := userMapping[un]; ok {
			return true
		}
	}
	return false
}

func setSession(w http.ResponseWriter, un string) {
	cookie := http.Cookie{
		Name:  "SID",
		Value: "test",
	}
	http.SetCookie(w, &cookie)
	sessions[cookie.Value] = un
}

func getUser(r *http.Request) (user, bool) {
	cookie, err := r.Cookie("SID")
	if err == nil {
		if un, ok := sessions[cookie.Value]; ok {
			if user, ok := userMapping[un]; ok {
				return user, true
			}
		}
	}
	return user{}, false
}
