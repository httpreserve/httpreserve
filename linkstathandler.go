package main

import (
	"os"
	"fmt"
	"net/http"
	"strings"
)

// Internal function used to finalize a struct to be used
// for reporting in the app whether our query has been a 
// successful one or not...
func makeLinkStats(ls LinkStats) {

	if !ls.ProtocolError {
		//fmt.Printf("%+v\n\n", ls)

		iaurllatest, err := GetPotentialUrlLatest(ls.Link.String())
		if err != nil {
			return
		}

		iaurlearliest, err := GetPotentialUrlEarliest(ls.Link.String())
		if err != nil {
			return
		}

      srlatest := CreateSimpleRequest(HEAD, iaurllatest, false, "")
		latestIA, err := httpFromSimpleRequest(srlatest)
		if err != nil {
			return
		}

		//first test for existence of an internet archive copy
		if latestIA.ResponseCode == http.StatusNotFound {
			fmt.Fprintln(os.Stderr, "[Internet Archive]", ERR_NO_IA)
		} else {
		   srearliest := CreateSimpleRequest(HEAD, iaurlearliest, false, "")
			earliestIA, err := httpFromSimpleRequest(srearliest)
			if err != nil {
				return
			}
			fmt.Println(latestIA.Link, latestIA.ResponseCode, latestIA.ResponseText)
			fmt.Println(earliestIA.Link, earliestIA.ResponseCode, earliestIA.ResponseText)

			fmt.Println()

			iaLinkData := earliestIA.header.Get("Link")
			iaLinkInfo := strings.Split(iaLinkData, ";")

			for _, m := range iaLinkInfo {
				fmt.Println(strings.Trim(m, " "))
			}

			//fmt.Printf("\n\n%s\n\n", iaLinkInfo)
			

			//fmt.Printf("\n%+v\n\n", earliestIA.header)

			//for _, m := range earliestIA.header {
			//	fmt.Println(m)
			//}
		}

		fmt.Println("---")
	} 
	// Else we can't add more information
}
