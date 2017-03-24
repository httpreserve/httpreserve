package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const requestedURL = "url"

// For debug, we have this function here just in case we need
// to take a look at our request headers...
func prettyRequest(w http.ResponseWriter, r *http.Request) {
	req, _ := httputil.DumpRequest(r, false)
	fmt.Fprintf(w, "%s", req)
	return
}

// Primary handler of all POST or GET requests to httpreserve
// pretty simple eh?!
func handleHttpreserve(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "httpreserve analysis:\n")
	switch r.Method {
	case http.MethodGet:
		lookup, _ := url.ParseQuery(r.URL.RawQuery)
		fmt.Fprintln(w, GenerateLinkStats(lookup[requestedURL][0]))
		return
	case http.MethodPost:
		r.ParseForm()
		fmt.Fprintln(w, GenerateLinkStats(r.Form.Get(requestedURL)))
		return
	}
}
