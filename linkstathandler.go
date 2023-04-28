package httpreserve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/httpreserve/phantomjsscreenshot"
	"github.com/httpreserve/simplerequest"
	"github.com/httpreserve/wayback"
	"github.com/pkg/errors"
)

var areyouthere = true

func init() {
	areyouthere = phantomjsscreenshot.Hello()
	if !areyouthere {
		// screenshot service isn't available
	}
}

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

// formatISODate returns a formatted date for a given value.
func formatISODate(inputDate string) string {
	const fourteenDigitDateFormat = "20060102150405"
	outputDate, _ := time.Parse(fourteenDigitDateFormat, inputDate)
	return fmt.Sprintf("%s", outputDate)
}

// getISODates make the 14-digit IA saved date a little more human
// readable and thus meaningful to users.
func getISODates(earliest string, latest string) (string, string) {

	// Compile a regex to match just 14-digit numeric values.
	regexPattern, _ := regexp.Compile("\\d{14}")

	datetime := regexPattern.FindString(earliest)

	if datetime == "" {
		return "", ""
	}

	earliest = formatISODate(datetime)

	// Latest date is implicit if there is an earliest.
	latest = formatISODate(regexPattern.FindString(latest))

	return earliest, latest
}

// Internal function used to finalize a struct to be used
// for reporting in the app whether our query has been a
// successful one or not...
func makeLinkStats(ls LinkStats, err error, encoded bool) (LinkStats, error) {
	ls.AnalysisVersionText = VersionText()
	ls.AnalysisVersionNumber = VersionNumber()
	ls.SimpleRequestVersion = simplerequest.Version()
	wb, err := wayback.GetWaybackData(ls.Link, VersionText())
	// else process the response and error...
	if wb.AlreadyWayback == nil {
		if wb.NotInWayback == false {
			ls.InternetArchiveLinkLatest = wb.LatestWayback
			earliest, latest := getISODates(wb.EarliestWayback, wb.LatestWayback)
			// There is only one snapshot and that will appear in the latest field.
			if !strings.Contains(wb.EarliestWayback, "web.archive.org/save") {
				ls.InternetArchiveLinkEarliest = wb.EarliestWayback
				ls.InternetArchiveEarliestDate = earliest
			}
			ls.InternetArchiveLatestDate = latest
			if !strings.Contains(latest, "web.archive.org/save") {
				var robustDateEarly, robustDateLate string
				// Add Robust links here.
				if !encoded {
					robustDateEarly, robustDateLate = getRobust(ls.Link, wb.EarliestWayback, wb.LatestWayback)
				} else {
					robustDateEarly, robustDateLate = getRobustEncoded(ls.Link, wb.EarliestWayback, wb.LatestWayback)
				}
				if robustDateEarly != "" {
					// There is only one snapshot and that will appear in the latest field.
					if !strings.Contains(wb.EarliestWayback, "web.archive.org/save") {
						ls.RobustLinkEarliest = robustDateEarly
					}
					ls.RobustLinkLatest = robustDateLate
				}
			}
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

	// attach a url for folks to save to wayback...
	ls.InternetArchiveSaveLink = wb.WaybackSaveURL

	// finally, add a screenshot to our LinkStats struct
	ls.ScreenShot = addScreenshot(ls)

	// how long did it take to process this record...
	ls = addTime(ls)

	return ls, nil
}

func addTime(ls LinkStats) LinkStats {
	elapsedTime = time.Since(startTime)
	ls.StatsCreationTime = elapsedTime.String()
	return ls
}

func addScreenshot(ls LinkStats) string {
	var link string
	var err error
	if snapshot == true && areyouthere == true {
		if ls.ResponseCode == 0 || ls.ResponseCode > 400 {
			link = ResponseIncorrect
			return link
		}
		link, err = phantomjsscreenshot.GrabScreenshot(ls.Link, 100, 100)
		if err != nil {
			if strings.Contains(link, phantomjsscreenshot.EncodingField) {
				//good chance we still have an image to decode
				return link
			}
			return GenerateSnapshotErr + " " + err.Error()
		}
	} else {
		link = SnapshotNotEnabled
	}
	return link
}

// Format our output to be useful to external callers
func makeLinkStatJSON(ls LinkStats) (string, error) {
	var buf = new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(ls)
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, buf.Bytes(), "", "   ")
	return string(prettyJSON.Bytes()), nil
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
