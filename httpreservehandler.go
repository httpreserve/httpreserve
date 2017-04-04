package httpreserve

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const requestedURL = "url"
const requestedFname = "filename"

const errParsingQuery = "error parsing query sent via GET"
const errNoURL = "no url specified, or too many"
const errMultiFname = "no filename, or more than one filename specified, setting to \"\""

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
		lookup, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			fmt.Fprintf(w, "%s\n", errParsingQuery)
		}

		var link string
		var fname string

		if val, ok := lookup[requestedURL]; ok {
			if len(val) > 0 && len(val) < 2 {
				link = val[0]
			}
		}

		if link == "" {
			fmt.Fprintf(w, "%s %s\n", errParsingQuery, errNoURL)
			return
		}

		if val, ok := lookup[requestedFname]; ok {
			if len(val) > 0 && len(val) < 2 {
				fname = val[0]
			}
		}

		if fname == "" {
			log.Printf("%s %s", r.RemoteAddr, errMultiFname)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, retrieveLinkStats(link, fname))
		return
	case http.MethodPost:
		r.ParseForm()
		link := r.Form.Get(requestedURL)
		fname := r.Form.Get(requestedFname)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, retrieveLinkStats(link, fname))
		return
	}
}

func retrieveLinkStats(link string, fname string) string {
	ls, _ := GenerateLinkStats(link, fname)
	js := MakeLinkStatsJSON(ls)
	return js
}
