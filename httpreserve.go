package httpreserve

import (
	"net/url"
	"github.com/pkg/errors"
)

// Httpreserves primary handler for different protocols
func testConnection(requrl string) (LinkStats, error) {
	var ls LinkStats
	var err error

	req, err := url.Parse(requrl)
	if err != nil {
		return ls, errors.Wrap(err, "url parse failed")
	}

	switch req.Scheme {
	case "ftp":
		/*_, err := handleftp(request)
		if err != nil {
			panic(err)
		}*/
	case "http":
		fallthrough
	case "https":
		ls, err = httpFromSimpleRequest(defaultSimpleRequest(req))
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
	ls.link = req
	ls.Link = req.String()
	return ls, nil
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
func GenerateLinkStats(link string) (LinkStats, error) {
	ls, err := linkStat(link)
	if err != nil {
		ls, _ = manageLinkStatErrors(ls, err)
		//TODO: consider what to do with manageLinkStatErrors here...
	}
	// Positive or negative result, populate LS structure
	ls, err = makeLinkStats(ls, err)
	if err != nil {
		if err.Error() == errorNoIALink {
			// we can ignore this, not a fatal error
		} else {
			return ls, err
		}
	}
	return ls, nil
}

// MakeLinkStatsJSON will output a LinkStats struct as 
// a JSON object to be used in our applications...
func MakeLinkStatsJSON(ls LinkStats) string {
	js, _ := makeLinkStatJSON(ls)
	return js
}

