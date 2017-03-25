package httpreserve

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const requestedURL = "url"

// A default value 
var defaultServerMethod = http.MethodPost

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
	switch r.Method {
	case http.MethodGet:
		lookup, _ := url.ParseQuery(r.URL.RawQuery)
		query := lookup[requestedURL][0]
		fmt.Fprintln(w, retrieveLinkStats(query))
		return
	case http.MethodPost:
		r.ParseForm()
		query := r.Form.Get(requestedURL)
		fmt.Fprintln(w, retrieveLinkStats(query))
		return
	}
}

func retrieveLinkStats(query string) string {
	ls, _ := GenerateLinkStats(query)
	js := MakeLinkStatsJSON(ls)
	return js
}
