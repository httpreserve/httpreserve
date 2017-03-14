package main

import (
	"github.com/pkg/errors"
	"encoding/base64"
	"crypto/tls"
	"net/http"
	"net/url"
)

func returnProxyClient(req *http.Request) (*http.Client, error) {

	c := &http.Client{}

	auth := "spencero:***"
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