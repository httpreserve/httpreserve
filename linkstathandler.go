package main

import (
	"strings"
	"net/http"
	"github.com/pkg/errors"
)

// Internal function used to finalize a struct to be used
// for reporting in the app whether our query has been a 
// successful one or not...
func makeLinkStats(ls LinkStats) (LinkStats, error) {

	if !ls.ProtocolError {

		iaurlearliest, err := GetPotentialUrlEarliest(ls.Link.String())
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

// Add the Internet Archive response codes to our structure
// for analysis outside of the package.
func addResponses(ls LinkStats, ia LinkStats) LinkStats {
	ls.InternetArchiveResponseCode = ia.ResponseCode
	ls.InternetArchiveResponseText = ia.ResponseText
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