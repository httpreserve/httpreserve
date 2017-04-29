package httpreserve

import (
	"encoding/json"
	"github.com/httpreserve/phantomjsscreenshot"
	"github.com/httpreserve/simplerequest"
	"github.com/httpreserve/wayback"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// GetLinkStatsHeader allows us to do some debug on the information
// returned from the server. First it mocks a response, and
// then adds some of our own information to it to enable
// DumpResponse prettyprint. We will consider its use in future
// As two pretty printed responses have been added to the struct.
func GetLinkStatsHeader(ls LinkStats) string {
	var r = http.Response{}
	r.StatusCode = ls.statuscode
	r.Status = ls.status
	r.Header = *ls.header
	req, _ := httputil.DumpResponse(&r, false)
	return string(req)
}

// GetLinkStatsURL returns the originally parsed URL
// as was sent to the server for a response.
func GetLinkStatsURL(ls LinkStats) *url.URL {
	return ls.link
}

// GetPrettyRequest returns the original request
// but pretty printed..
func GetPrettyRequest(ls LinkStats) string {
	return ls.prettyRequest
}

// GetPrettyResponse returns the original response
// but pretty printed.
func GetPrettyResponse(ls LinkStats) string {
	return ls.prettyResponse
}

// Internal function used to finalize a struct to be used
// for reporting in the app whether our query has been a
// successful one or not...
func makeLinkStats(ls LinkStats, err error) (LinkStats, error) {

	ls.AnalysisVersionText = VersionText()
	ls.AnalysisVersionNumber = VersionNumber()
	ls.SimpleRequestVersion = simplerequest.Version()

	wb, err := wayback.GetWaybackData(ls.Link, VersionText())
	// else process the response and error...
	if wb.AlreadyWayback == nil {
		if wb.NotInWayback == false {
			ls.InternetArchiveLinkEarliest = wb.EarliestWayback
			ls.InternetArchiveLinkLatest = wb.LatestWayback
		} else {
			ls.Archived = false
		}

		ls.InternetArchiveResponseCode = wb.ResponseCode

		ls.InternetArchiveResponseText = wb.ResponseText
		if err != nil {
			ls.InternetArchiveResponseText = err.Error()
		}

		// plus a bit more to understand if the link is archived
		if !(ls.InternetArchiveResponseCode == http.StatusNotFound || ls.InternetArchiveResponseCode == 0) {
			ls.Archived = true
		}
	}

	ls.InternetArchiveSaveLink = wb.WaybackSaveURL

	//finally, add a screenshot to our LinkStats struct
	ls.ScreenShot = addScreenshot(ls)

	return ls, nil
}

func addScreenshot(ls LinkStats) string {
	if snapshot {
		link, err := phantomjsscreenshot.GrabScreenshot(ls.Link)
		if err != nil {
			if strings.Contains(link, phantomjsscreenshot.EncodingField) {
				//good chance we still have an image to decode
				return link
			}
			return err.Error()
		}
	}
	return snapshotmessage
}

// Format our output to be useful to external callers
func makeLinkStatJSON(ls LinkStats) (string, error) {
	js, err := json.MarshalIndent(ls, "", "   ")
	if err != nil {
		return "", err
	}
	return string(js), nil
}

// Add important errors to LinkStats structure for us to
// work with as and when we need to, example, zero protocol
func manageLinkStatErrors(ls LinkStats, err error) (LinkStats, error) {
	//report error in some way...
	ls.Error = true
	switch err.Error() {
	case errorBlankProtocol:
		ls.ErrorMessage = errorBlankProtocol
	case errorUnknownProtocol:
		ls.ErrorMessage = errorUnknownProtocol
	default:
		return ls, errors.Wrap(err, "LinkStat failed")
	}
	return ls, nil
}
