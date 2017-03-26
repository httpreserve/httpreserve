package httpreserve

import (
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// At least for testing we're going to be doing a limited range
// of things with our requests. Create a default object to make that
// easier for us.
func defaultSimpleRequest(reqURL *url.URL) SimpleRequest {
	return CreateSimpleRequest(httpHEAD, reqURL, useProxy, httpBYTERANGE)
}

// CreateSimpleRequest is a mechanism to make a suitable
// http request header to find some information out about
// a web resouse.
// We want to make handlehttp more useable so let's wrap
// as much as we can up front and see if that's possible
// recommended setting for byterange is to maintain the default
// but the potential to set it manually here is possible
// If byterange is left "" then default range will be used.
func CreateSimpleRequest(method string, reqURL *url.URL, proxy bool, byterange string) SimpleRequest {
	var sr SimpleRequest
	sr.Method = method
	sr.ReqURL = reqURL
	sr.Proxy = proxy
	if byterange == "" {
		sr.ByteRange = httpBYTERANGE
	} else {
		sr.ByteRange = byterange
	}
	return sr
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
		return ls, errors.Wrap(err, "client request failed")
	}

	// A mechanism for users to debug their code using Response headers
	re, _ := httputil.DumpResponse(resp, false)
	ls.prettyResponse = string(re)

	ls.ResponseCode = resp.StatusCode
	ls.ResponseText = http.StatusText(resp.StatusCode)
	ls.link = reqURL
	ls.Link = reqURL.String()

	// For debug record pertinent packet details...
	ls.statuscode = resp.StatusCode
	ls.status = resp.Status
	ls.header = &resp.Header

	// Do we have to do NT lan Manager negotiation...
	if checkNTLM(resp, reqURL) {
		resp.Body.Close()
		return ls, errors.New(errorNTLM)
	}

	// once we've closed the body we can't do anything else
	// with the packet content...
	resp.Body.Close()

	return ls, nil
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
