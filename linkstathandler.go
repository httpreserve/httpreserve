package main

import "fmt"

// Internal function used to finalize a struct to be used
// for reporting in the app whether our query has been a 
// successful one or not...
func makeLinkStats(ls LinkStats) {
	fmt.Printf("%+v\n", ls)
}
