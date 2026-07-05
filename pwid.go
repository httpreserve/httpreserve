package httpreserve

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// formatPWIDDate returns a Robust formatted date for a given value.
func formatPWIDDate(inputDate string) string {
	const fourteenDigitDateFormat = "20060102150405"
	outputDate, _ := time.Parse(fourteenDigitDateFormat, inputDate)
	return fmt.Sprintf("%s", outputDate.Format("2006-01-02T15:04:05Z"))
}

// getPWIDDate returns a PWID formatted date to the caller.
func getPWIDDate(date string) string {

	// Compile a regex to match just 14-digit numeric values.
	regexPattern, _ := regexp.Compile("\\d{14}")
	formattedDate := formatPWIDDate(regexPattern.FindString(date))
	return formattedDate
}

// getPWID attempts to construct a persistent web identifier for a given
// resource.
//
// The logic needs some work. We will look for text/html and text/plain
// and application/xhtml+xml to begin with.
//
// Page: urn:pwid:archive.org:2016-01-22T10:08:23Z:page:https://www.dr.dk
// Part: urn:pwid:archive.org:2022-12-12T17:14:47Z:part:http://id.kb.dk/pwid/PWID.ppsm
func getPWID(link string, contentType string, date string) string {

	pageMimes := []string{
		// MIMETypes associated with entire web-pages, not web-parts.
		"text/plain",
		"text/html",
		"application/xhtml+xml",
	}

	// Top 20 or so TDLs to help disambiguate page versus part.
	tlds := []string{
		".app",
		".at",
		".au",
		".be",
		".berlin",
		".biz",
		".br",
		".ca",
		".ch",
		".cn",
		".com",
		".co.uk",
		".de",
		".es",
		".eu",
		".fr",
		".gov",
		".in",
		".io",
		".info",
		".it",
		".net",
		".nl",
		".online",
		".org",
		".org",
		".pl",
		".rocks",
		".ru",
		".shop",
		".store",
		".tk",
		".tv",
		".uk",
		".us",
		".xyz",
	}

	// NB. Keep in mind that this will change with a different memento.
	const urn string = "urn:pwid:archive.org"
	pwidDate := getPWIDDate(date)

	if contentType != "" {
		for _, val := range pageMimes {
			if strings.Contains(contentType, val) {
				if !strings.HasSuffix(link, "/") {
					link = fmt.Sprintf("%s/", link)
				}
				return fmt.Sprintf("%s:%s:page:%s", urn, pwidDate, link)
			}
		}
	}

	// We haven't got a content type. Use a heuristic to try and
	// derive more info.
	for _, value := range tlds {
		if strings.HasSuffix(link, value) {
			return fmt.Sprintf("%s:%s:page:%s", urn, pwidDate, link)
		}
	}

	if strings.HasSuffix(link, "/") ||
		strings.HasSuffix(link, ".htm") ||
		strings.HasSuffix(link, ".html") ||
		strings.HasSuffix(link, ".xhtml") ||
		strings.HasSuffix(link, ".md") ||
		strings.HasSuffix(link, ".php") ||
		strings.HasSuffix(link, ".aspx") {
		return fmt.Sprintf("%s:%s:page:%s", urn, pwidDate, link)
	}

	// Finally, assume page part.
	return fmt.Sprintf("%s:%s:part:%s", urn, pwidDate, link)
}
