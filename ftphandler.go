package httpreserve

import (
	"fmt"
	"github.com/dutchcoders/goftp"
)

// handleftp for eventually returning some kind of FTP 
// response if we're asked to...
// TODO: lots to still implement here...
func handleftp(request string) (LinkStats, error) {
	//var ftp *goftp.FTP

	// For debug messages: goftp.ConnectDbg("ftp.server.com:21")
	//"ftp.server.com:21"

	request = request + ":21"
	request = "ftp.exponentialdecay.co.uk:21"

	fmt.Println(request)
	//goftp.ConnectDbg(
	//TODO: parse debug messages...
	if _, err := goftp.ConnectDbg(request); err != nil {
		checkError(err)
	}

	var ls LinkStats
	return ls, nil
}

// Helper function to manage errors...
func checkError(err error) {
	if err != nil {
		fmt.Println("[FTP]", err.Error())
		return
	}
}
