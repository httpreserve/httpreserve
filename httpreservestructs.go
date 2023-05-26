package httpreserve

import (
	"net/http"
	"net/url"
)

// LinkStats Table structure to be returned from our
// requests Can be fairly liberal in its expansion
type LinkStats struct {
	FileName                    string `json:"FileName,omitempty"` // If a filename is provided
	AnalysisVersionNumber       string
	AnalysisVersionText         string
	SimpleRequestVersion        string
	Link                        string
	Title                       string
	ContentType                 string
	ResponseCode                int
	ResponseText                string
	SourceURL                   string // URL requested by the caller
	ScreenShot                  string // HREF to screenshot
	InternetArchiveLinkEarliest string `json:"InternetArchiveLinkEarliest,omitempty"`
	InternetArchiveEarliestDate string `json:"InternetArchiveEarliestDate,omitempty"`
	InternetArchiveLinkLatest   string
	InternetArchiveLatestDate   string `json:"InternetArchiveLatestDate,omitempty"`
	InternetArchiveSaveLink     string // Link to use to save from the Internet
	InternetArchiveResponseCode int
	InternetArchiveResponseText string
	RobustLinkEarliest          string `json:"RobustLinkEarliest,omitempty"` // A robust hyperlink snippet linking to a live url and a memento version
	RobustLinkLatest            string `json:"RobustLinkLatest,omitempty"`   // A robust hyperlink snippet linking to a live url and a memento version
	PWID                        string `json:PWID,omitempty`                 // Persistent Web Identifier DRAFT URN standard from Denmark.
	Archived                    bool   // Has the Internet Archive saved the page or not?
	Error                       bool
	ErrorMessage                string
	StatsCreationTime           string

	// For debug
	status     string
	statuscode int
	header     *http.Header
	link       *url.URL
	tld        string

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
