package routes

import (
	"log"
	"net/http"
	"regexp"
)

var routeMap = make(map[string]map[string]http.HandlerFunc)

func init() {
	for _, route := range Routes {
		if method, ok := route[0].(string); ok {
			// if check if method exists in http methods
			if routeStr, ok := route[1].(string); ok {
				if handleFunc, ok := route[2].(func(http.ResponseWriter, *http.Request)); ok {
					if routeMap[method] == nil {
						routeMap[method] = map[string]http.HandlerFunc{routeStr: handleFunc}
					} else {
						routeMap[method][routeStr] = handleFunc
					}
				} else {
					log.Printf("%T is not of type func(w http.ResponseWriter, r *http.Request)\n", handleFunc)
				}
			}
		}
	}
}

func myfunc(w http.ResponseWriter, r *http.Request) {
	h := http.NotFound
	for k, v := range routeMap[r.Method] {
		// fmt.Printf(" %v ==> %v\n", regexp.MustCompile("^"+k+"$").MatchString(r.URL.String()), regexp.MustCompile("^"+k+"/$").MatchString(r.URL.String()))
		if regexp.MustCompile("^"+k+"$").MatchString(r.URL.String()) || regexp.MustCompile("^"+k+"/$").MatchString(r.URL.String()) {
			h = v
			break
		}
	}
	h(w, r)
}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", myfunc)
	log.Println("Listening....")
	http.ListenAndServe(":3000", mux)
}
