package main 

import (
	"strings"
	"net/url"
	"net/http"
	"github.com/pkg/errors"
)

const IA_ROOT = "http://web.archive.org"

func GetSavedURL(resp http.Response) (*url.URL, error) {
	loc := resp.Header["Content-Location"]
	u, err := url.Parse(IA_ROOT + strings.Join(loc, ""))
	if err != nil {
		return &url.URL{}, errors.Wrap(err, "creation of URL from http response failed.")
	}
	return u, nil
}