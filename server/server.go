package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/justinas/alice"
)

func httpreserve(w http.ResponseWriter, r *http.Request) { 
	fmt.Fprintf(w, "Some information: %s\n", "gah!")
	//http.ServeHTTP(w, r)
}

func fourohfour(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintln(w, "This is not the primary entry point.")
	//http.ServeHTTP(w, r)
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requested %s, method %s", r.RemoteAddr, r.URL, r.Method)
		h.ServeHTTP(w, r)
	})
}

type headerSetter struct {
	key, val string
	handler  http.Handler
}

func (hs headerSetter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(hs.key, hs.val)
	hs.handler.ServeHTTP(w, r)
}

func newHeaderSetter(key, val string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return headerSetter{key, val, h}
	}
}

//reference: https://cryptic.io/go-http/
//https://github.com/justinas/alice
func main() {
	h := http.NewServeMux()

	h.HandleFunc("/httpreserve", httpreserve)
	h.HandleFunc("/", fourohfour) 

	//Middleware chain to handle various generic HTTP functions
	middleware_chain := alice.New(
		newHeaderSetter("Server", "exponentialDK-httpreserve/0.0.0"),	// USERAGENT IN MAIN PACKAGE
		logger,
	).Then(h)

	err := http.ListenAndServe(":2040", middleware_chain)
	log.Fatal(err)
}
