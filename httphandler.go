package httpreserve

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/httpreserve/simplerequest"
	"github.com/pkg/errors"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// use a HEAD request to calibrate our own request to the server
// return an integer if successful, and false for no request error
// return zero if not successful, and true for error...
func configureRequest(sr simplerequest.SimpleRequest) (string, bool) {

	// first utilize the simplerequest we already have
	srHead := sr
	srHead.Method = simplerequest.HEAD

	// try and ilicit a response from the server
	resp, err := srHead.Do()
	if err != nil {
		return "", true
	}

	cl := resp.Header.Get("Content-Length")
	if cl == "" {
		return "", false
	}

	i, err := strconv.Atoi(cl)
	if err != nil {
		return "", false
	}

	// understand how we want to work with zeros length...
	if i == 0 {
		return "", false
	}

	// now return the minimum value we can retrieve from server
	i = min(i, 500)

	return strconv.Itoa(i), false
}

// HTTPFromSimpleRequest is another mechanism we can use to
// retrieve some basic information out from a web resource.
// Call handlehttp from a SimpleRequest object instead
// of calling function directly...
func HTTPFromSimpleRequest(sr simplerequest.SimpleRequest, sourceURL string) (LinkStats, error) {

	// identify our agent, and then configure requesr...
	sr.Agent(VersionText())

	// configure our values...
	byterange, e := configureRequest(sr)

	//set some values for the simplerequest...
	if e {
		sr.Timeout(5) // fail quick
	} else {
		sr.Timeout(10) // take our time if more potential
	}

	if byterange != "" {
		sr.Byterange(byterange)
	}

	//retrieve our link stats...
	return getLinkStats(sr, sourceURL)
}

// retrieveTLD returns what should be the top-level domain for a given
// URL.
func retrieveTLD(link *url.URL) string {
	hostSplit := strings.Split(link.Hostname(), ".")
	if len(hostSplit) < 2 {
		// Initial thoughts about this is that there isn't much we can
		// do here, but also, we're unlikely to find a host that doesn't
		// split into at least two parts - this is probably some other
		// error that we've picked up earlier.
		//
		// NB. defensive.
		//
		return link.Hostname()
	}
	return strings.Join(hostSplit[len(hostSplit)-2:], ".")
}

// Handle HTTP functions of the calling application.
func getLinkStats(req simplerequest.SimpleRequest, sourceURL string) (LinkStats, error) {
	var ls LinkStats

	// populate linkstats asap...
	ls.link = req.URL
	ls.Link = req.URL.String()
	ls.tld = retrieveTLD(req.URL)

	if sourceURL == "" {
		sourceURL = ls.Link
	}

	// make sure if we get a url shortener we handle it on its merit...
	req.NoRedirect(true)

	sr, err := req.Do()
	if err != nil {
		if strings.Contains(err.Error(), "lookup") &&
			strings.Contains(err.Error(), "no such host") {
			ls.ResponseText = "error: client request failed: no such host"
		} else if strings.Contains(err.Error(), "i/o timeout") {
			ls.ResponseText = "error: client request failed: i/o timeout"
		} else if strings.Contains(err.Error(), "no route to host") {
			ls.ResponseText = "error: client request failed: no route to host"
		} else {
			ls.ResponseText = "client request failed: " + err.Error()
		}
		// return and only continue to proces responses that there
		// was no error for...
		return ls, err
	}

	// We may have a url shortening service...
	if sr.Location != nil && (sr.StatusCode >= 300 && sr.StatusCode < 400) {
		// If the top-level domain isn't in the new location, we will
		// treat the request as being a request from a URL-shortener,
		// that means, the TLD is elsewhere.
		if !strings.Contains(sr.Location.String(), ls.tld) {
			// If the tld and the returned location is different, then
			// we should check if it is a shortner to follow. If it
			// isn't a shortner, then we should request the URL as-is.
			//
			// NB. we don't need these nested ifs. We can ask straight
			// out if it's a shortener. But I quite like them for now
			// for readability and extensibility and I', not worried
			// about performance.
			//
			if _, ok := shorteners[ls.tld]; ok {
				return HTTPFromSimpleRequest(simplerequest.Default(sr.Location), sourceURL)
			}
		}
	}

	// start adding to our LinkStat struct as soon as possible
	ls.link = req.URL
	ls.Link = req.URL.String()

	//Get our pretty printed output for debug etc.
	ls.prettyRequest = sr.PrettyRequest
	ls.prettyResponse = sr.PrettyResponse

	// Response Codes...
	ls.ResponseCode = sr.StatusCode
	ls.ResponseText = sr.StatusText

	// Populate LS Title and Content-Type
	ls.ContentType = sr.GetHeader("Content-Type")

	// Look at the payload to see if we can retrieve title...
	ls.Title = getTitle(string(sr.Data), ls.ContentType)

	// For debug record pertinent packet details...
	ls.header = &sr.Header

	// Commit the original request to the stat here.
	ls.SourceURL = sourceURL

	// Do we have to do NT lan Manager negotiation...
	if checkNTLM(sr) {
		return ls, errors.New(errorNTLM)
	}

	return ls, nil
}

// GetTitle is a way to add more useful metadata to our LinkStats
// structure by way of checking for link drift. Where the page we're
// expecting is one thing but it has become another.
func getTitle(body string, contentType string) string {
	if !strings.Contains(contentType, "text/html") {
		return ""
	}
	body = strings.ToLower(body)
	t1string := "<title>"
	t1 := strings.Index(body, t1string)
	t2 := strings.Index(body, "</title>")
	if (t1 != -1 && t2 != -1) && t2 > t1+len(t1string) {
		return body[t1+len(t1string) : t2] //index plus length of search string
	}
	return ""
}

// Network back-ends like here at Archives New Zealand use NTLM
// authentication as a secondary proxy that applications have to
// jump through. NTLM stands for NT Lan Management. If we receive
// a cue to have to do NTLM authentication then we need to jump
// through those hoops. We begin that process here.
func checkNTLM(sr simplerequest.SimpleResponse) bool {
	if sr.StatusCode == 407 {
		if strings.Join(sr.Header[authNTLM], " ") == flagNTLM {
			// we have to do the NTLM DANCE here...
			// https://github.com/exponential-decay/httpreserve/issues/1
			return true
		}
	}
	return false
}
