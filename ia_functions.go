package main 

import (
	"strings"
	"net/url"
	"net/http"
	"github.com/pkg/errors"
)

const IA_ROOT = "http://web.archive.org"

// Create a URL that we can test for a 404 error or 200 OK. 
// The URL if it works can be used to display to the user for
// QA. The URL if it fails, can be used to prompt the user to
// save the URL as it is found today. A motivation, even if there
// is no saved IA record, to save copy today, even if it is a 404
// is that the earliest date we can pin on a broken link the 
// better we can satisfy outselves in future that we did all we can.
func GetPotentialUrlLatest() {

}

// There may be benefit to returning the earliest possible record
// available in the internet archive. We can make it easier by
// using this function here. 
func GetPotentialUrlEarliest() {

}

// Construct the url to return to either the IA earliest or latest
// IA get functions and return...
func constructUrl(iadate string) string {

	return ""
}

// Utilize the methods across the package to submit a URL to the 
// internet archive to retrieve a saved URL that we can use.
func SubmitToInternetarchive() {

}

// We've constructed the URL to save ours in the Internet Archive
// We've submitted the URL via the IA REST API and we've receieved
// a 200 OK. In the response will be a partial SLUG that takes us
// to our newly archived record. 
func GetSavedURL(resp http.Response) (*url.URL, error) {
	loc := resp.Header["Content-Location"]
	u, err := url.Parse(IA_ROOT + strings.Join(loc, ""))
	if err != nil {
		return &url.URL{}, errors.Wrap(err, "creation of URL from http response failed.")
	}
	return u, nil
}

