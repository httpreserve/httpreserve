package main

import (
	"os"
	"fmt"
	"strings"
	"net/http"
	"encoding/json"	
	"github.com/pkg/errors"
)

// Internal function used to finalize a struct to be used
// for reporting in the app whether our query has been a 
// successful one or not...
func makeLinkStats(ls LinkStats, err error) (LinkStats, error) {

	if !ls.ProtocolError {

		iaurlearliest, err := GetPotentialUrlEarliest(ls.Link)
		if err != nil {
			return ls, errors.Wrap(err, "IA url creation failed")
		}

		srearliest := CreateSimpleRequest(HEAD, iaurlearliest, false, "")
		earliestIA, err := httpFromSimpleRequest(srearliest)
		if err != nil {
			return ls, errors.Wrap(err, "IA http request failed")
		}

		// Add out Internet Archive Response Code to ours...
		ls = addResponses(ls, earliestIA)
		ls = addSaveUrl(ls)

		// First test for existence of an internet archive copy
		if earliestIA.ResponseCode == http.StatusNotFound {

			if earliestIA.header.Get("Link") == "" {
				return ls, errors.New(ERR_NO_IA)
			}
		} 

		// Else, continue to retrieve IA links
		iaLinkData := earliestIA.header.Get("Link")
		iaLinkInfo := strings.Split(iaLinkData, ", <")

		var legacy_collection = make(map[string]string)

		for _, lnk := range iaLinkInfo {
			trimmedlink := strings.Trim(lnk, " ")
			trimmedlink = strings.Replace(trimmedlink, ">;", ";", 1)		// fix chevrons
			for  _, rel := range IA_REL_LIST {
				if strings.Contains(trimmedlink, rel) {
					legacy_collection[rel] = trimmedlink
					break
				}
			}
		}

		// We've some internet archive links that we can use
		if len(legacy_collection) > 0 {
			ls = populateIALinks(ls, legacy_collection)
		}
	} 
	return ls, nil
}

// Format our output to be useful to external callers
func makeLinkStatJson(ls LinkStats) (string, error) {
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
	case ERR_BLANK_PROTOCOL:
		ls.ProtocolErrorMessage = ERR_BLANK_PROTOCOL 
	case ERR_UNKNOWN_PROTOCOL:
		ls.ProtocolErrorMessage = ERR_UNKNOWN_PROTOCOL
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
func addSaveUrl(ls LinkStats) LinkStats {
	ls.InternetArchiveSaveLink = makeSaveUrl(ls.Link)
	return ls
}

// We are interested in the earliest link available in the
// Internet Archive and the latest link available, return the strings here.
func populateIALinks(ls LinkStats, legacy_collection map[string]string) LinkStats {
	for rel, lnk := range legacy_collection {
		switch rel {
		// first two cases give us the earliest IA link available
		case REL_FIRST:
			fallthrough
		case REL_FIRST_LAST:
			ls.InternetArchiveLinkEarliest = getWWW(lnk)
			break
		// second two cases give us the latest IA link available
		case REL_LAST:
			fallthrough
		case REL_NEXT_LAST:
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
		if strings.Contains(www, IA_ROOT) {
			return www
		}
	}
	return ""
}