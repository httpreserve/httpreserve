package httpreserve

import (
	"fmt"
	"testing"
)

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
