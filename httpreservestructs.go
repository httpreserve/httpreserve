package main

import (
	"net/http"
	"net/url"
)

// HTTP request methods that are useful to us
const GET = http.MethodGet
const HEAD = http.MethodHead

// User-agent to identify code being run
const USERAGENT = "exponentialDK-httpreserve/0.0.0"

// Default byte-range for initial requests
const BYTERANGE = "bytes=0-0"

// Default proxy value we might set on compilation
const USE_PROXY = false

// SimpleRequest structure to be turned into a
// HTTP request proper in code.
type SimpleRequest struct {
	Method								string
	ReqUrl								*url.URL
	Proxy									bool
	ByteRange							string
}

// Table structure to be returned from our requests
// Can be fairly liberal in its expansion
type LinkStats struct {
	FileName 							string	// If a filename is provided
	Link 									string
	ResponseCode 						int
	ResponseText 						string
	ScreenShot 							string 	// HREF to screenshot
	InternetArchiveLinkLatest		string
	InternetArchiveLinkEarliest	string	// Earliest link in Internet Archive
	InternetArchiveSaveLink			string	// Link to use to save from the Internet
	InternetArchiveResponseCode 	int
	InternetArchiveResponseText 	string
	Archived								bool		// Has the Internet Archive saved the page or not?
	ProtocolError						bool
	ProtocolErrorMessage				string

	// For debug
	header *http.Header
	link   *url.URL	
}

// NTLM (NT Lan Management) Consts
// For when we can configure this code to run against NTLM
var AUTH_CODE = 407
var NTLM_AUTH = "Proxy-Authenticate"
var NTLM_FLAG = "Negotiate NTLM"
var ERR_NTLM  = "Requires NTLM Negotiation"

