package main

import "net/http"

type LinkStats struct {
	FileName 							string
	Link 									string
	ResponseCode 						int
	ResponseText 						string
	ScreenShot 							string //href to screenshot
	InternetArchiveLink 				string
	InternetArchiveResponseCode 	int
	InternetArchiveResponseText 	string
	ArchiveNow							bool

	//for debug
	header *http.Header
}
