package httpreserve

import (
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const iaRoot = "http://web.archive.org"
const iaBeta = "http://web-beta.archive.org"

const iaSRoot = "https://web.archive.org"
const iaSBeta = "https://web-beta.archive.org"

const iaSave = "/save/" //e.g. https://web.archive.org/save/http://www.bbc.com/news
const iaWeb = "/web/"   //e.g. http://web.archive.org/web/20161104020243/http://exponentialdecayxxxx.co.uk/#
const iaRel = "rel="

// Memento returns various relationship attributes
// These are the ones observed so far in this work.
// Rather than separating the attributes, use whole string.
const relFirst = "rel=\"first memento\""          // syn: first, at least two
const relNext = "rel=\"next memento\""            // syn: at least three
const relLast = "rel=\"last memento\""            // syn: last, at least three
const relFirstLast = "rel=\"first last memento\"" // syn: only
const relNextLast = "rel=\"next last memento\""   // syn: second, and last
const relPrevLast = "rel=\"prev memento\""        // syn: at least three
const relPrevFirst = "rel=\"prev first memento\"" // syn: previous, and first, only two

// List of items to check against when parsing header attributes
var iaRelList = [...]string{relFirst, relNext, relLast, relFirstLast,
	relNextLast, relPrevLast, relPrevFirst}

//Explanation: https://andrey.nering.com.br/2015/how-to-format-date-and-time-with-go-lang/
//Golang Date Formatter: http://fuckinggodateformat.com/
const datelayout = "20060102150405"

// GetPotentialURLLatest is used to create a URL that we can test for a 404 
// error or 200 OK. The URL if it works can be used to display to 
// the user for QA. The URL if it fails, can be used to prompt the 
// user to save the URL as it is found today. A motivation, even if 
// there is no saved IA record, to save copy today, even if it is a 404
// is that the earliest date we can pin on a broken link the
// better we can satisfy outselves in future that we did all we can.
// Example URI we need to create looks like this:
// web.archive.org/web/{date}/url-to-lookup
// {date} == "20161104020243" == "YYYYMMDDHHMMSS" == %Y%m%d%k%M%S
func GetPotentialURLLatest(archiveurl string) (*url.URL, error) {
	latestDate := time.Now().Format(datelayout)
	return constructURL(latestDate, archiveurl)
}

// GetPotentialURLEarliest is used to returning the 
// earliest possible record available in the internet archive. We 
// can make it easier by using this function here.
// Example URI we need to create looks like this:
// web.archive.org/web/{date}/url-to-lookup
func GetPotentialURLEarliest(archiveurl string) (*url.URL, error) {
	oldestDate := time.Date(1900, time.August, 31, 23, 13, 0, 0, time.Local).Format(datelayout)
	return constructURL(oldestDate, archiveurl)
}

// Construct the url to return to either the IA earliest or latest
// IA get functions and return...
func constructURL(iadate string, archiveurl string) (*url.URL, error) {
	newurl, err := url.Parse(iaRoot + iaWeb + iadate + "/" + archiveurl)
	if err != nil {
		return newurl, errors.Wrap(err, "internet archive url creation failed")
	}
	return newurl, nil
}

// MakeSaveURL is used to create a URL that will enable us to 
// submit it to the Internet Archive SaveNow function
func MakeSaveURL(link string) string {
	//e.g. https://web.archive.org/save/http://www.bbc.com/news
	return iaRoot + iaSave + link
}

// SubmitToInternetArchive will handle the request and response to
// and from the Internet Archive for a URL that we wish to save as
// part of this initiative. 
func SubmitToInternetArchive() {

}

// GetSavedURL will help us to retrieve the URL returned by the 
// Internet Archive when we've sent a request to the SaveNow function.
// We've constructed the URL to save ours in the Internet Archive
// We've submitted the URL via the IA REST API and we've receieved
// a 200 OK. In the response will be a partial SLUG that takes us
// to our newly archived record.
func GetSavedURL(resp http.Response) (*url.URL, error) {
	loc := resp.Header["Content-Location"]
	u, err := url.Parse(iaRoot + strings.Join(loc, ""))
	if err != nil {
		return &url.URL{}, errors.Wrap(err, "creation of URL from http response failed.")
	}
	return u, nil
}

// Test the URL to make sure it's not already an internet archive link
func isIA(link string) bool {
	if strings.Contains(link, iaRoot) || strings.Contains(link, iaBeta) ||
		strings.Contains(link, iaSRoot) || strings.Contains(link, iaSBeta) {
		return true
	}
	return false
}
