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

var AUTH_CODE = 407
var NTLM_AUTH = "Proxy-Authenticate"
var NTLM_FLAG = "Negotiate NTLM"
var ERR_NTLM = "Requires NTLM Negotiation"

//proxy help:
//https://jannewmarch.gitbooks.io/network-programming-with-go-golang-/content/http/proxy_handling.html
//another example:
//http://stackoverflow.com/questions/40817784/access-https-via-http-proxy-with-basic-authentication
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
	resp.Body.Close()
y


	//fmt.Println(resp)
	//fmt.Printf("%+v\n", resp)

	if checkNTLM(resp) {
		return ls, errors.New(ERR_NTLM)
	}



	return ls, nil
}

func checkNTLM(resp *http.Response) bool {
	if resp.StatusCode == 407 {
		fmt.Print("x ", strings.Join(resp.Header[NTLM_AUTH], ""), " x\n")
		if strings.Join(resp.Header[NTLM_AUTH], " ") == NTLM_FLAG {
			//NTLM DANCE
			//https://msdn.microsoft.com/en-us/library/dd925287(v=office.12).aspx
			fmt.Println("NTLM Dance Here")
			return true
		}
	}
	return false
}