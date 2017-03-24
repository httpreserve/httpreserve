package main

import (
	"fmt"
	"testing"
)

//sample working perma.cc
//https://perma.cc/T8U2-994F
//https://perma.cc/48VC-ZS62
//https://perma.cc/9265-T4NB

//broken perma.cc (should be a 404)
//https://perma.cc/48VC-ZS61

func TestIALinkNowDate(t *testing.T) {
	archiveurl, _ := GetPotentialURLLatest("http://www.bbc.co.uk/news")
	fmt.Println(archiveurl)
}

func TestIALinkEarliestDate(t *testing.T) {
	archiveurl, _ := GetPotentialURLEarliest("http://www.bbc.co.uk/news")
	fmt.Println(archiveurl)
}

func TestSavedURL(t *testing.T) {
	u, _ := GetSavedURL(generateInternetArchiveSaveMock())
	fmt.Println(u)
}

func TestPlaceHolder(t *testing.T) {

}
