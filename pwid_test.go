package httpreserve

import (
	"testing"
)

type pwidTest struct {
	url         string
	contentType string
	date        string
	result      string
}

var pwidTests = []pwidTest{
	{
		"http://example.com",
		"application/html",
		"20170908143755",
		"urn:pwid:archive.org:2017-09-08T14:37:55Z:page:http://example.com",
	},
	{
		"http://example.com/file.jpg",
		"applicatio/html",
		"20170908143755",
		"urn:pwid:archive.org:2017-09-08T14:37:55Z:part:http://example.com/file.jpg",
	},
	{
		"http://example.com/",
		"applicatio/html",
		"19700101010101",
		"urn:pwid:archive.org:1970-01-01T01:01:01Z:page:http://example.com/",
	},
	{
		"http://example.com/",
		"applicatio/html",
		"20260705000100",
		"urn:pwid:archive.org:2026-07-05T00:01:00Z:page:http://example.com/",
	},
}

func TestGetPWID(t *testing.T) {
	/* Example:

	urn:pwid:archive.org:2026-07-03T03:19:07Z:page:https://example.com/"
	*/

	for _, v := range pwidTests {
		pwid := getPWID(v.url, v.contentType, v.date)
		if pwid != v.result {
			t.Fatalf("res '%s' does not equal expected: '%s'", pwid, v.result)
		}
	}
}
