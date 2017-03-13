package main

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
}
