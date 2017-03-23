package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/justinas/alice"
)

// Primary handler for httpreserve requests
func httpreserve(w http.ResponseWriter, r *http.Request) { 
	fmt.Fprintf(w, "Some information: %s\n", "gah!")
}

// Return a 404: TODO: May discard in favour of more friendly
// response for the user...
func fourohfour(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintln(w, "This is not the primary entry point.")
}

// Logger middleware to return information to stderr we're
// interested in...
func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requested %s, method %s", r.RemoteAddr, r.URL, r.Method)
		h.ServeHTTP(w, r)
	})
}

// Part of our Handler Adapter methods
// TODO: learn more about to document further
type headerSetter struct {
	key, val string
	handler  http.Handler
}

// Part of middleware layer to create default header responses
func (hs headerSetter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(hs.key, hs.val)
	hs.handler.ServeHTTP(w, r)
}

// Set default headers for any single response from httpreserve
func newHeaderSetter(key, val string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return headerSetter{key, val, h}
	}
}

// Configure our default server mechanism for httpreserve
func configureDefault() http.Handler {
	h := http.NewServeMux()

	h.HandleFunc("/httpreserve", httpreserve)
	h.HandleFunc("/", fourohfour) 

	// Middleware chain to handle various generic HTTP functions
	// TODO: Learn what other middleware we may need...
	middleware_chain := alice.New(
		newHeaderSetter("Server", "exponentialDK-httpreserve/0.0.0"),	// USERAGENT IN MAIN PACKAGE
		logger,
	).Then(h)

	return middleware_chain
}

// Standup a default server for the httpreserve resolver
// service to be queried by our other apps.
func DefaultServer(port string) {
	mw := configureDefault()
	err := http.ListenAndServe(":" + port, mw)
	log.Fatal(err)
}

// References contributing to this code...
// https://cryptic.io/go-http/
// https://github.com/justinas/alice
func main() {
	DefaultServer("2040")
}
