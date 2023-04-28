package httpreserve

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// formatRobustDate returns a Robust formatted date for a given value.
func formatRobustDate(inputDate string) string {
	const fourteenDigitDateFormat = "20060102150405"
	outputDate, _ := time.Parse(fourteenDigitDateFormat, inputDate)
	return fmt.Sprintf("%s", outputDate.Format("2006-01-02"))
}

// getRobustDates returns a simplified date format according to the
// RobustLinks specification.
//
// TODO: Some of this code is repeated from getISODates. We should be
// able to grab the dates early and then process without having to do
// the regext work or checking twice.
func getRobustDates(earliest string, latest string) (string, string) {

	// Compile a regex to match just 14-digit numeric values.
	regexPattern, _ := regexp.Compile("\\d{14}")

	datetime := regexPattern.FindString(earliest)

	if datetime == "" {
		return "", ""
	}

	earliest = formatRobustDate(datetime)

	// Latest date is implicit if there is an earliest.
	latest = formatRobustDate(regexPattern.FindString(latest))

	return earliest, latest
}

// makeReplacements will make the neccesary string replacements in our
// robust links.
func makeReplacements(html string, origin string, ia string, date string) string {

	replaceDate := "{{ date }}"
	replaceOrigin := "{{ origin }}"
	replaceIA := "{{ ia }}"

	html = strings.Replace(html, replaceDate, date, 1)
	html = strings.Replace(html, replaceOrigin, origin, 1)
	html = strings.Replace(html, replaceIA, ia, 1)

	return html
}

// getRobust will return earliest and latest Robust links for the given
// resource, which looks something like follows:
//
//		/*
//			"RobustLinkEarliest": "
//				<a href=\"http://web.archive.org/web/20020120142510/http://example.com/\"
//				 data-originalurl=\"http://example.com/\"
//				 data-versiondate=\"2002-01-20\">
//				 HTTPreserve Robust Link - simply replace this text!</a>",
//		*/
//
//	 NB. HTML fragment validator: https://appdevtools.com/html-validator
func getRobust(origin string, earliest string, latest string) (string, string) {

	earliestDate, latestDate := getRobustDates(earliest, latest)

	if earliestDate == "" {
		return "", ""
	}

	early := "<a href='{{ ia }}'" +
		"\x20data-originalurl='{{ origin }}'" +
		"\x20data-versiondate='{{ date }}'>" +
		"HTTPreserve Robust Link - simply replace this text!!</a>"

	late := "<a href='{{ ia }}'" +
		"\x20data-originalurl='{{ origin }}'" +
		"\x20data-versiondate='{{ date }}'>" +
		"HTTPreserve Robust Link - simply replace this text!!</a>"

	early = makeReplacements(early, origin, earliest, earliestDate)
	late = makeReplacements(late, origin, latest, latestDate)

	return early, late
}

func getRobustEncoded(origin string, earliest string, latest string) (string, string) {

	earliestDate, latestDate := getRobustDates(earliest, latest)

	if earliestDate == "" {
		return "", ""
	}

	early := "&lt;a href='{{ ia }}'" +
		"\x20data-originalurl='{{ origin }}'" +
		"\x20data-versiondate='{{ date }}'&gt;" +
		"HTTPreserve Robust Link - simply replace this text!!&lt;/a&gt;"

	late := "&lt;a href='{{ ia }}'" +
		"\x20data-originalurl='{{ origin }}'" +
		"\x20data-versiondate='{{ date }}'&gt;" +
		"HTTPreserve Robust Link - simply replace this text!!&lt;/a&gt;"

	early = makeReplacements(early, origin, earliest, earliestDate)
	late = makeReplacements(late, origin, latest, latestDate)

	return early, late
}
