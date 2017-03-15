package main 

import (
	"os"
	"fmt"
	"github.com/dutchcoders/goftp"	
)

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
		panic(err)
	}

	var ls LinkStats
	return ls, nil
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}