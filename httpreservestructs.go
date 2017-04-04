package httpreserve

import (
	"net/http"
	"net/url"
)

// HTTP request methods that are useful to us
const httpGET = http.MethodGet
const httpHEAD = http.MethodHead

// Default byte-range for initial requests
const httpBYTERANGE = "bytes=0-0"

// Default proxy value we might set on compilation
const useProxy = false

// SimpleRequest structure to be turned into a
// HTTP request proper in code.
type SimpleRequest struct {
	Method    string
	ReqURL    *url.URL
	Proxy     bool
	ByteRange string
}

// LinkStats Table structure to be returned from our
// requests Can be fairly liberal in its expansion
type LinkStats struct {
	FileName                    string // If a filename is provided
	AnalysisVersionNumber       string
	AnalysisVersionText         string
	Link                        string
	Title                       string
	ContentType                 string
	ResponseCode                int
	ResponseText                string
	ScreenShot                  string // HREF to screenshot
	InternetArchiveLinkLatest   string
	InternetArchiveLinkEarliest string // Earliest link in Internet Archive
	InternetArchiveSaveLink     string // Link to use to save from the Internet
	InternetArchiveResponseCode int
	InternetArchiveResponseText string
	Archived                    bool // Has the Internet Archive saved the page or not?
	ProtocolError               bool
	ProtocolErrorMessage        string

	// For debug
	status     string
	statuscode int
	header     *http.Header
	link       *url.URL

	// Pretty Debug
	prettyRequest  string
	prettyResponse string
}

// NTLM (NT Lan Management) Consts
// For when we can configure this code to run against NTLM
var authCode = 407
var authNTLM = "Proxy-Authenticate"
var flagNTLM = "Negotiate NTLM"
var errorNTLM = "Requires NTLM Negotiation"
