package httpreserve

import (
	"github.com/httpreserve/simplerequest"
	"github.com/httpreserve/wayback"
	"github.com/pkg/errors"
	"net/url"
	"time"
)

var startTime time.Time
var elapsedTime time.Duration

// handleRequest is HTTPreserve's primary handler for different protocols.
func handleRequest(reqURL string) (LinkStats, error) {
	var ls LinkStats
	var err error

	req, err := url.Parse(reqURL)
	if err != nil {
		return ls, errors.Wrap(err, "url parse failed")
	}

	// Option to handle FTP if we choose, but deleted from code
	// at present, before git: 26a91577bc7bb8d29187169755e01caf730b2f14
	switch req.Scheme {
	case "http":
		fallthrough
	case "https":
		startTime = time.Now() // Record processing time.
		ls, err := HTTPFromSimpleRequest(simplerequest.Default(req), "")
		if err != nil {
			return ls, errors.Wrap(err, "handleRequest() failed")
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
	ls, err := handleRequest(url)
	return ls, err
}

// generateLinkStats performs the setup work prior to returning a LinkStat struct.
func generateLinkStats(link string) LinkStats {
	ls, err := linkStat(link)
	if err != nil {
		ls, _ = manageLinkStatErrors(ls, err)
		//TODO: consider what to do with manageLinkStatErrors here...
	}
	return ls
}

// handleLinkStatError is a helper for our generateLinkStats functions below.
func handleLinkStatError(ls LinkStats, err error) (LinkStats, error) {
	if err.Error() == wayback.ErrorNoIALink.Error() {
		// we can ignore this, not a fatal error.
		return ls, nil
	}
	return ls, err
}

// GenerateLinkStats is used to return a JSON object for a URL
// specified in link variable passed to the function.
func GenerateLinkStats(link string, fileName string, screenGrab bool) (LinkStats, error) {
	// Limit data being sent by disabling screenshots if selected.
	if screenGrab != true {
		snapshot = false
	}
	ls := generateLinkStats(link)
	// Positive or negative result, populate LS structure
	var err error
	ls, err = makeLinkStats(ls, err, false)
	if err != nil {
		_, err = handleLinkStatError(ls, err)
		if err != nil {
			return ls, err
		}
	}
	if fileName != "" {
		ls.FileName = fileName
	}
	return ls, nil
}

// GenerateLinkStatsEncoded encodes LinkStat JSON with HTML entities for display online,
// e.g. on HTTPreserve.info.
func GenerateLinkStatsEncoded(link string, fileName string, screenGrab bool) (LinkStats, error) {
	// Limit data being sent by disabling screenshots if selected.
	if screenGrab != true {
		snapshot = false
	}
	ls := generateLinkStats(link)
	// Positive or negative result, populate LS structure
	var err error
	ls, err = makeLinkStats(ls, err, true)
	if err != nil {
		_, err = handleLinkStatError(ls, err)
		if err != nil {
			return ls, err
		}
	}
	if fileName != "" {
		ls.FileName = fileName
	}
	return ls, nil
}

// MakeLinkStatsJSON will output a LinkStats struct as
// a JSON object to be used in our applications...
func MakeLinkStatsJSON(ls LinkStats) string {
	js, _ := makeLinkStatJSON(ls)
	return js
}
