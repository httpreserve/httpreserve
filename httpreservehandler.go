package httpreserve

import (
	"fmt"
	"github.com/httpreserve/simplerequest"
	"github.com/httpreserve/wayback"
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

// Use this function to retrieve all the args sent to the handler
func getLinkFname(w http.ResponseWriter, r *http.Request) (string, string, string) {

	var link string
	var fname string

	switch r.Method {
	case http.MethodGet:
		lookup, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			return "", "", errParsingQuery
		}

		if val, ok := lookup[requestedURL]; ok {
			if len(val) > 0 && len(val) < 2 {
				link = val[0]
			}
		}

		if link == "" {
			return "", "", errNoURL
		}

		if val, ok := lookup[requestedFname]; ok {
			if len(val) > 0 && len(val) < 2 {
				fname = val[0]
			}
		}

	case http.MethodPost:
		r.ParseForm()
		link = r.Form.Get(requestedURL)
		fname = r.Form.Get(requestedFname)
	}

	return link, fname, ""
}

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
	// get our variable values
	link, fname, e := getLinkFname(w, r)
	if e != "" {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, e)
		return
	}

	// push json to client
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, retrieveLinkStats(link, fname))
	return
}

// test for a url shortening service
func checkshort(link string) string {
	u, _ := url.Parse(link)
	sr := simplerequest.Default(u)
	sr.NoRedirect(true)
	resp, _ := sr.Do()
	if resp.StatusCode == 301 && resp.Location != nil {
		return resp.Location.String()
	}
	return ""
}

// submit link to internet archive
func handleSubmitToInternetArchive(w http.ResponseWriter, r *http.Request) {

	// push json to client
	w.Header().Set("Content-Type", "application/json")

	var ls LinkStats

	// get our variable values
	link, fname, e := getLinkFname(w, r)
	if e != "" {
		ls.Error = true
		ls.ErrorMessage = e
		fmt.Fprintln(w, MakeLinkStatsJSON(ls))
		return
	}

	unshort := checkshort(link)
	if unshort != "" {
		link = unshort
	}

	// else continue to submit to internet archive
	_, err := wayback.SubmitToInternetArchive(link, VersionText())
	if err != nil {
		ls.FileName = fname
		ls.Link = link
		ls.Error = true
		ls.ErrorMessage = "saving link to archive didn't work, " + err.Error()
		fmt.Fprintln(w, MakeLinkStatsJSON(ls))
		return
	}

	fmt.Fprintln(w, retrieveLinkStats(link, fname))
	return
}

// retrieve linkstats from httpreserve
func retrieveLinkStats(link string, fname string) string {
	ls, _ := GenerateLinkStatsEncoded(link, fname, true)
	js := MakeLinkStatsJSON(ls)
	return js
}
