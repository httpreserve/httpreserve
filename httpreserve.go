package httpreserve

import (
	"github.com/httpreserve/simplerequest"
	"github.com/httpreserve/wayback"
	"github.com/pkg/errors"
	"net/url"
	"time"
)

var starttime time.Time
var elapsedtime time.Duration

// Httpreserves primary handler for different protocols
func testConnection(requrl string) (LinkStats, error) {
	var ls LinkStats
	var err error

	req, err := url.Parse(requrl)
	if err != nil {
		return ls, errors.Wrap(err, "url parse failed")
	}

	// Option to handle FTP if we choose, but deleted from code
	// at present, before git: 26a91577bc7bb8d29187169755e01caf730b2f14
	switch req.Scheme {
	case "http":
		fallthrough
	case "https":
		// processing time
		starttime = time.Now()

		ls, err := HTTPFromSimpleRequest(simplerequest.Default(req))
		if err != nil {
			return ls, errors.Wrap(err, "handlehttp() failed")
		}
		return ls, nil
	case "":
		ls.link = req
		ls.Link = req.String()
		return ls, errors.New(errorBlankProtocol)
	default:
		ls.link = req
		ls.Link = req.String()
		return ls, errors.Wrap(errors.New(errorUnknownProtocol), req.Scheme)
	}
}

// Function that wraps testconnection, we may refactor this at
// some point as it doesn't do a lot in its own right.
func linkStat(url string) (LinkStats, error) {
	var ls LinkStats
	ls, err := testConnection(url)
	return ls, err
}

// GenerateLinkStats is used to return a JSON object for a URL
// specified in link variable passed to the function.
func GenerateLinkStats(link string, fname string, screengrab bool) (LinkStats, error) {

	// set global variable to help folks limit data sent by screenshots
	if screengrab != true {
		snapshot = false
	}

	ls, err := linkStat(link)
	if err != nil {
		ls, _ = manageLinkStatErrors(ls, err)
		//TODO: consider what to do with manageLinkStatErrors here...
	}
	// Positive or negative result, populate LS structure
	ls, err = makeLinkStats(ls, err)
	if err != nil {
		if err.Error() == wayback.ErrorNoIALink.Error() { // TODO: may be able to remove
			// we can ignore this, not a fatal error
		} else {
			return ls, err
		}
	}
	if fname != "" {
		ls.FileName = fname
	}
	return ls, nil
}

// MakeLinkStatsJSON will output a LinkStats struct as
// a JSON object to be used in our applications...
func MakeLinkStatsJSON(ls LinkStats) string {
	js, _ := makeLinkStatJSON(ls)
	return js
}
