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

const GET = http.MethodGet
const HEAD = http.MethodHead

var AUTH_CODE = 407
var NTLM_AUTH = "Proxy-Authenticate"
var NTLM_FLAG = "Negotiate NTLM"
var ERR_NTLM  = "Requires NTLM Negotiation"

// Handle HTTP functions of the calling application. If we need to use
// a proxy then set the flag, if not, then don't. 
func handlehttp(request string, proxflag bool) (LinkStats, error) {

	var ls LinkStats
	var client = &http.Client{}

	linkurl, err := url.Parse(request)
	if err != nil {
		return ls, errors.Wrap(err, "parse request url failed")
	}

	req, err := http.NewRequest(HEAD, linkurl.String(), nil)
	if err != nil {
		return ls, errors.Wrap(err, "request generation failed")
	}
   req.Header.Add("User-Agent", USERAGENT)
   req.Header.Add("Range", BYTERANGE) 
   req.Header.Add("proxy-Connection", "Keep-Alive")

	if proxflag {
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

	if checkNTLM(resp, request) {
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
func checkNTLM(resp *http.Response, request string) bool {
	if resp.StatusCode == 407 {
		if strings.Join(resp.Header[NTLM_AUTH], " ") == NTLM_FLAG {			
			// we have to do the NTLM DANCE here...
			// https://github.com/exponential-decay/httpreserve/issues/1
			return true
		}
	}
	return false
}