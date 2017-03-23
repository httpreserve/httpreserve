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

// 404 response handler for all non supported function
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "Sorry, this is not a supported function for this application.")
}

// Return a 404: TODO: May discard in favour of more friendly
// response for the user...
func indexhandler(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(404)
	if r.URL.String() != "/" {
		NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodOptions:
		handleOptions(w, r)
		return
	case http.MethodHead:
		fallthrough
	case http.MethodGet:
		fmt.Fprintln(w, "Deliver static page?")	
		return	
	default:
		fmt.Fprintln(w, r.Method + " is unsupported from root.")
		return
	}
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

// Handle the return of favicon when requested by client
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "ico/favicon.ico")
}

// Configure our default server mechanism for httpreserve
func configureDefault() http.Handler {
	h := http.NewServeMux()

	h.HandleFunc("/favicon.ico", faviconHandler)
	h.HandleFunc("/httpreserve", httpreserve)
	h.HandleFunc("/", indexhandler) 

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
