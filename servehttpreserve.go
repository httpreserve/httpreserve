package main

import (
	"fmt"
	"net/url"
	"net/http"
	"net/http/httputil"
)

const REQUESTED_URL = "url"

func handleHttpreserve(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Return server httpreserve analysis here.\n")
	req, _ := httputil.DumpRequest(r, false)
	fmt.Fprintf(w, "%s", req)

	switch r.Method {
	case http.MethodGet:
		lookup, _ := url.ParseQuery(r.URL.RawQuery)
		fmt.Fprintln(w, lookup[REQUESTED_URL][0])
		return
	case http.MethodPost:
		r.ParseForm()
		fmt.Fprintln(w, "\n\n" + r.Form.Get("url"))
		return
	}
}
