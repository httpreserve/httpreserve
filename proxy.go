package main

import (
	"github.com/pkg/errors"
	"encoding/base64"
	"crypto/tls"
	"net/http"
	"net/url"
)

// Links that may help with proxy configuration for our application
// https://jannewmarch.gitbooks.io/network-programming-with-go-golang-/content/http/proxy_handling.html
// http://stackoverflow.com/questions/40817784/access-https-via-http-proxy-with-basic-authentication

func returnProxyClient(req *http.Request) (*http.Client, error) {

	c := &http.Client{}

	auth := "dia\\spencero:****"
	proxy := "https://wlgproxy:8080"

	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return c, errors.Wrap(err, "parse proxy url failed")
	}

	// encode the auth
	basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Proxy-Authorisation", basic)

	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	transport.ProxyConnectHeader = req.Header

	c = &http.Client{Transport: transport}	

	return c, nil
}