package main

import (
	"fmt"
	"net/url"
	"net/http"
	"net/http/httputil"
)

func handleHttpreserve(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Return server httpreserve analysis here.\n")
	req, _ := httputil.DumpRequest(r, false)
	fmt.Fprintf(w, "%s", req)

	switch r.Method {
	case http.MethodGet:
		m, _ := url.ParseQuery(r.URL.String())
		fmt.Fprintf(w, "\n\n%+v\n\n", m)
		//fmt.Println(m["k"][0])
		return
	case http.MethodPost:
		r.ParseForm()
		fmt.Fprintln(w, "\n\n" + r.Form.Get("url"))
		return
	}
}
