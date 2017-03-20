package main 

import (
	"os"
	"fmt"
	"strings"
	"net/url"
	"net/http"
	"net/http/httputil"
	"github.com/pkg/errors"
)

// At least for testing we're going to be doing a limited range
// of things with our requests. Create a default object to make that
// easier for us.
func DefaultSimpleRequest(requrl string) (SimpleRequest, error) {
	sr, err := CreateSimpleRequest(HEAD, requrl, USE_PROXY, BYTERANGE)
	if err != nil {
		return sr, errors.Wrap(err, "default simpleUrl url assignment failed")
	}
	return sr, nil
}

// We want to make handlehttp more useable so let's wrap
// as much as we can up front and see if that's possible
// recommended setting for byterange is to maintain the default
// but the potential to set it manually here is possible
func CreateSimpleRequest(method string, newUrl string, proxy bool, byterange string) (SimpleRequest, error) {
	var sr SimpleRequest
	sr.Method = method
	reqUrl, err := url.Parse(newUrl)
	if err != nil {
		return sr, errors.Wrap(err, "simpleUrl url assignment failed")
	}
	sr.ReqUrl = reqUrl
	sr.Proxy = proxy
	if byterange == "" {
		sr.ByteRange = BYTERANGE
	} else {
		sr.ByteRange = byterange
	}
	return sr, nil
}

// Call handlehttp from a SimpleRequest object instead
// of calling function directly...
func httpFromSimpleRequest(sr SimpleRequest) (LinkStats, error) {
	ls, err := handlehttp(sr.Method, sr.ReqUrl, sr.Proxy, sr.ByteRange)
	return ls, err	
}

// Handle HTTP functions of the calling application. If we need to use
// a proxy then set the flag, if not, then don't. 
func handlehttp(method string, requrl *url.URL, proxy bool, byterange string) (LinkStats, error) {

	var ls LinkStats
	var client = &http.Client{}

	req, err := http.NewRequest(method, requrl.String(), nil)
	if err != nil {
		return ls, errors.Wrap(err, "request generation failed")
	}
   req.Header.Add("User-Agent", USERAGENT)
   req.Header.Add("Range", byterange) 
   req.Header.Add("proxy-Connection", "Keep-Alive")

	if proxy {
		client, err = returnProxyClient(req)	
		if err != nil {
			return ls, errors.Wrap(err, "proxy header creation failed")
		}
	} 

	dump, _ := httputil.DumpRequest(req, false)
	fmt.Println("Request header:")
	fmt.Fprintln(os.Stdout, string(dump))

	resp, err := client.Do(req)
	if err != nil {
		return ls, errors.Wrap(err, "client request failed")
	}

	ls.ResponseCode = resp.StatusCode
	ls.ResponseText = http.StatusText(resp.StatusCode)

	if checkNTLM(resp, requrl) {
		resp.Body.Close()
		return ls, errors.New(ERR_NTLM)
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
func checkNTLM(resp *http.Response, reqUrl *url.URL) bool {
	if resp.StatusCode == 407 {
		if strings.Join(resp.Header[NTLM_AUTH], " ") == NTLM_FLAG {			
			// we have to do the NTLM DANCE here...
			// https://github.com/exponential-decay/httpreserve/issues/1
			return true
		}
	}
	return false
}
