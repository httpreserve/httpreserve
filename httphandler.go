package httpreserve

import (
	"github.com/httpreserve/simplerequest"
	"github.com/pkg/errors"
	"strings"
)

// HTTPFromSimpleRequest is another mechanism we can use to
// retrieve some basic information out from a web resource.
// Call handlehttp from a SimpleRequest object instead
// of calling function directly...
func HTTPFromSimpleRequest(sr simplerequest.SimpleRequest) (LinkStats, error) {
	
	//set some values for the simplerequest...
	sr.Timeout(10)
	sr.Agent(VersionText())
	sr.Byterange("500")

	//retrieve our link stats...
	ls, err := getLinkStats(sr)
	return ls, err
}

// Handle HTTP functions of the calling application.
func getLinkStats(req simplerequest.SimpleRequest) (LinkStats, error) {
	// We're going to have some data to work with so lets 
	// start populating our LinkStats struct
	var ls LinkStats
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
