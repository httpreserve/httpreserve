package httpreserve

import (
	"os"
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"net/http/httputil"
	"github.com/pkg/errors"
)

// GetLSHeader allows us to do some debug on the information
// returned from the server. First it mocks a response, and 
// then adds some of our info to it to enable DumpResponse prettyprint
func GetLSHeader(ls LinkStats) string {
	var r = http.Response{}
	r.StatusCode = ls.statuscode
	r.Status = ls.status
	r.Header = *ls.header 
	req, _ := httputil.DumpResponse(&r, false)
	return string(req)
}

// Internal function used to finalize a struct to be used
// for reporting in the app whether our query has been a
// successful one or not...
func makeLinkStats(ls LinkStats, err error) (LinkStats, error) {

	if !ls.ProtocolError {

		iaURLearliest, err := GetPotentialURLEarliest(ls.Link)
		if err != nil {
			return ls, errors.Wrap(err, "IA url creation failed")
		}

		if !isIA(ls.Link) {
			isEarliest := createSimpleRequest(httpHEAD, iaURLearliest, false, "")
			earliestIA, err := httpFromSimpleRequest(isEarliest)
			if err != nil {
				return ls, errors.Wrap(err, "IA http request failed")
			}
			// Add out Internet Archive Response Code to ours...
			ls = addResponses(ls, earliestIA)

			// First test for existence of an internet archive copy
			if earliestIA.ResponseCode == http.StatusNotFound {
				if earliestIA.header.Get("Link") == "" {
					return ls, errors.New(errorNoIALink)
				}
			}

			// Else, continue to retrieve IA links
			iaLinkData := earliestIA.header.Get("Link")
			iaLinkInfo := strings.Split(iaLinkData, ", <")

			var legacyCollection = make(map[string]string)

			for _, lnk := range iaLinkInfo {
				trimmedlink := strings.Trim(lnk, " ")
				trimmedlink = strings.Replace(trimmedlink, ">;", ";", 1) // fix chevrons
				for _, rel := range iaRelList {
					if strings.Contains(trimmedlink, rel) {
						legacyCollection[rel] = trimmedlink
						break
					}
				}
			}

			// We've some internet archive links that we can use
			if len(legacyCollection) > 0 {
				ls = populateIALinks(ls, legacyCollection)
			}

			ls = addSaveURL(ls)
		} else {
			ls.InternetArchiveSaveLink = ErrorIAExists
		}
	}
	return ls, nil
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
	ls.ProtocolError = true
	switch err.Error() {
	case errorBlankProtocol:
		ls.ProtocolErrorMessage = errorBlankProtocol
	case errorUnknownProtocol:
		ls.ProtocolErrorMessage = errorUnknownProtocol
	default:
		fmt.Fprintln(os.Stderr, "[LinkStat Fail]", errors.Wrap(err, "LinkStat failed"))
	}
	return ls, nil
}

// Add the Internet Archive response codes to our structure
// for analysis outside of the package.
func addResponses(ls LinkStats, ia LinkStats) LinkStats {
	ls.InternetArchiveResponseCode = ia.ResponseCode
	ls.InternetArchiveResponseText = ia.ResponseText
	if ia.ResponseCode == http.StatusNotFound || ia.ResponseCode == 0 {
		ls.Archived = false
	} else {
		ls.Archived = true
	}
	return ls
}

// Add the Internet Archive save link to our linkstat struct
// to enable saving of the most up-to-date version of the resource
func addSaveURL(ls LinkStats) LinkStats {
	ls.InternetArchiveSaveLink = MakeSaveURL(ls.Link)
	return ls
}

// We are interested in the earliest link available in the
// Internet Archive and the latest link available, return the strings here.
func populateIALinks(ls LinkStats, legacyCollection map[string]string) LinkStats {
	for rel, lnk := range legacyCollection {
		switch rel {
		// first two cases give us the earliest IA link available
		case relFirst:
			fallthrough
		case relFirstLast:
			ls.InternetArchiveLinkEarliest = getWWW(lnk)
			break
		// second two cases give us the latest IA link available
		case relLast:
			fallthrough
		case relNextLast:
			ls.InternetArchiveLinkLatest = getWWW(lnk)
			break
		}
	}
	return ls
}

// Retrieve the IA www link that we've been passing about
// from the IA response header sent to us previously.
func getWWW(lnk string) string {
	lnksplit := strings.Split(lnk, "; ")
	for _, www := range lnksplit {
		if strings.Contains(www, iaRoot) {
			return www
		}
	}
	return ""
}
