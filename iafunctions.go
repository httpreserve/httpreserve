package main 

import (
	"time"
	"strings"
	"net/url"
	"net/http"
	"github.com/pkg/errors"
)

const IA_ROOT = "http://web.archive.org"
const IA_SAVE = "/save/"		//e.g. https://web.archive.org/save/http://www.bbc.com/news
const IA_WEB = "/web/"			//e.g. http://web.archive.org/web/20161104020243/http://exponentialdecayxxxx.co.uk/#
const IA_REL = "rel="

// Memento returns various relationship attributes
// These are the ones observed so far in this work.
// Rather than separating the attributes, use whole string.
const REL_FIRST = "rel=\"first memento\""					// syn: first, at least two
const REL_NEXT = "rel=\"next memento\""					// syn: at least three
const REL_LAST = "rel=\"last memento\""					// syn: last, at least three
const REL_FIRST_LAST = "rel=\"first last memento\""	// syn: only
const REL_NEXT_LAST = "rel=\"next last memento\""		// syn: second, and last
const REL_PREV_LAST = "rel=\"prev memento\""				// syn: at least three
const REL_PREV_FIRST = "rel=\"prev first memento\""	// syn: previous, and first, only two

// List of items to check against when parsing header attributes
var IA_REL_LIST = [...]string{REL_FIRST, REL_NEXT, REL_LAST, REL_FIRST_LAST, 
		REL_NEXT_LAST, REL_PREV_LAST, REL_PREV_FIRST}

//Explanation: https://andrey.nering.com.br/2015/how-to-format-date-and-time-with-go-lang/
//Golang Date Formatter: http://fuckinggodateformat.com/
const DATELAYOUT = "20060102150405"

// Create a URL that we can test for a 404 error or 200 OK. 
// The URL if it works can be used to display to the user for
// QA. The URL if it fails, can be used to prompt the user to
// save the URL as it is found today. A motivation, even if there
// is no saved IA record, to save copy today, even if it is a 404
// is that the earliest date we can pin on a broken link the 
// better we can satisfy outselves in future that we did all we can.
// Example URI we need to create looks like this:
// web.archive.org/web/{date}/url-to-lookup
// {date} == "20161104020243" == "YYYYMMDDHHMMSS" == %Y%m%d%k%M%S
func GetPotentialUrlLatest(archiveurl string) (*url.URL, error) {
	latestDate := time.Now().Format(DATELAYOUT)
	return constructUrl(latestDate, archiveurl)
}

// There may be benefit to returning the earliest possible record
// available in the internet archive. We can make it easier by
// using this function here. 
// Example URI we need to create looks like this:
// web.archive.org/web/{date}/url-to-lookup
func GetPotentialUrlEarliest(archiveurl string) (*url.URL, error) {
	oldestDate := time.Date(1900, time.August, 31, 23, 13, 0, 0, time.Local).Format(DATELAYOUT)
	return constructUrl(oldestDate, archiveurl)
}

// Construct the url to return to either the IA earliest or latest
// IA get functions and return...
func constructUrl(iadate string, archiveurl string) (*url.URL, error) {
	newurl, err := url.Parse(IA_ROOT + IA_WEB + iadate + "/" + archiveurl)
	if err != nil {
		return newurl, errors.Wrap(err, "internet archive url creation failed")
	}
	return newurl, nil
}

// Create a URL that will enable us to submit it to the Internet Archive
func makeSaveUrl(link string) string {
	//e.g. https://web.archive.org/save/http://www.bbc.com/news
	if strings.Contains(link, IA_ROOT) {
		return ERR_ALREADY_IA  						// validity of error? seems useful.
	}
	return IA_ROOT + IA_SAVE + link
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
