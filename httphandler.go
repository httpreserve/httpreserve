package httpreserve

import (
	//"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// At least for testing we're going to be doing a limited range
// of things with our requests. Create a default object to make that
// easier for us.
func defaultSimpleRequest(reqURL *url.URL) SimpleRequest {
	// we're not concerned about error here, as internally, we've
	// already parsed the URL which is the only source of potential
	// error in CreateSimpleRequest
	sr, _ := CreateSimpleRequest(httpGET, reqURL.String(), useProxy, httpBYTERANGE)
	return sr
}

// CreateSimpleRequest is a mechanism to make a suitable
// http request header to find some information out about
// a web resouse.
// We want to make handlehttp more useable so let's wrap
// as much as we can up front and see if that's possible
// recommended setting for byterange is to maintain the default
// but the potential to set it manually here is possible
// If byterange is left "" then default range will be used.
func CreateSimpleRequest(method string, reqURL string, proxy bool, byterange string) (SimpleRequest, error) {
	var sr SimpleRequest
	sr.Method = method
	req, err := url.Parse(reqURL)
	if err != nil {
		return sr, errors.Wrap(err, "url parse failed in CreateSimpleRequest")
	}
	sr.ReqURL = req
	sr.Proxy = proxy
	if byterange == "" {
		sr.ByteRange = httpBYTERANGE
	} else {
		sr.ByteRange = byterange
	}
	return sr, nil
}

// HTTPFromSimpleRequest is another mechanism we can use to
// retrieve some basic information out from a web resource.
// Call handlehttp from a SimpleRequest object instead
// of calling function directly...
func HTTPFromSimpleRequest(sr SimpleRequest) (LinkStats, error) {
	ls, err := handlehttp(sr.Method, sr.ReqURL, sr.Proxy, sr.ByteRange)
	return ls, err
}

// Handle HTTP functions of the calling application. If we need to use
// a proxy then set the flag, if not, then don't.
func handlehttp(method string, reqURL *url.URL, proxy bool, byterange string) (LinkStats, error) {

	var ls LinkStats
	var client = &http.Client{}

	req, err := http.NewRequest(method, reqURL.String(), nil)
	if err != nil {
		return ls, errors.Wrap(err, "request generation failed")
	}
	req.Header.Add("User-Agent", VersionText())
	req.Header.Add("Range", byterange)
	req.Header.Add("proxy-Connection", "Keep-Alive")

	// start adding to our LinkStat struct as soon as possible
	ls.link = reqURL
	ls.Link = reqURL.String()

	if proxy {
		client, err = returnProxyClient(req)
		if err != nil {
			return ls, errors.Wrap(err, "proxy header creation failed")
		}
	}

	// A mechanism for users to debug their code using Request headers
	rq, _ := httputil.DumpRequest(req, false)
	ls.prettyRequest = string(rq)

	resp, err := client.Do(req)
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

	// once we've closed the body we can't do anything else
	// with the packet content...
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return ls, errors.Wrap(err, "reading http response body")
	}

	// A mechanism for users to debug their code using Response headers
	re, _ := httputil.DumpResponse(resp, false)
	ls.prettyResponse = string(re)

	// Response Codes...
	ls.ResponseCode = resp.StatusCode
	ls.ResponseText = http.StatusText(resp.StatusCode)

	// Populate LS Title and Content-Type
	ls.ContentType = resp.Header.Get("Content-Type")
	ls.Title = getTitle(string(data), ls.ContentType)

	// For debug record pertinent packet details...
	ls.statuscode = resp.StatusCode
	ls.status = resp.Status
	ls.header = &resp.Header

	// Do we have to do NT lan Manager negotiation...
	if checkNTLM(resp, reqURL) {
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
func checkNTLM(resp *http.Response, reqURL *url.URL) bool {
	if resp.StatusCode == 407 {
		if strings.Join(resp.Header[authNTLM], " ") == flagNTLM {
			// we have to do the NTLM DANCE here...
			// https://github.com/exponential-decay/httpreserve/issues/1
			return true
		}
	}
	return false
}
