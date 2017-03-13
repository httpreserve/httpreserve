package main

import (
		"os"
      "net/http"
      "fmt"
		//"log"
		//"io/ioutil"
      "bufio"
		"net/url"
		"github.com/pkg/errors"
		"github.com/dutchcoders/goftp"
   )

//ftp: ftp://exponentialdecay.co.uk/
//http; http://exponentialdecay.co.uk
//https: https://github.com/exponential-decay

//hackable uri for archive.org, finds closest to date set, pre or post
//internet archive uri: http://web.archive.org/web/20161104020243/http://exponentialdecay.co.uk/#

const CONN_OKAY int8 = 0
const CONN_BAD int8 = 1

const GET = http.MethodGet
const HEAD = http.MethodHead

const USERAGENT = "@exponentialDK httpreserve"
const BYTERANGE = "bytes=0-0"

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

func handlehttp(request string) (LinkStats, error) {

	var ls LinkStats

	req, err := http.NewRequest(HEAD, request, nil) 
	if err != nil {
		return ls, errors.Wrap(err, "create request failed")
		//if unsupported protocol... correct? 
	}

   req.Header.Add("User-Agent", USERAGENT)
   req.Header.Add("Range", BYTERANGE) 

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ls, errors.Wrap(err, "client request failed")
	}

	ls.ResponseCode = resp.StatusCode
	ls.ResponseText = http.StatusText(resp.StatusCode)
	resp.Body.Close()

	return ls, nil

}

func testConnection (request string) (LinkStats, error) {

	var ls LinkStats

	u, err := url.Parse(request)
	if err != nil {
		return ls, errors.Wrap(err, "url parse failed")
	}

	fmt.Println(u.Scheme)
	if u.Scheme == "ftp" {
		_, err := handleftp(request)
		if err != nil {
			panic(err)
		}
	}

	return ls, nil
}

func LinkStat(url string) (LinkStats, error) {
	var a LinkStats
	a, err := testConnection(url)
	return a, err
}

//might be useful...
func stdiolooper() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "[x]", errors.Wrap(err, "read failed"))
	}
}

func looper() {

	file, err := os.Open("link-examples/linklist.txt") 
	if err != nil {
		fmt.Fprintln(os.Stderr, "[x]", errors.Wrap(err, "file open failed"))
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		//send throught to our function tog get stats...
		//fmt.Println(scanner.Text()) 
		ls, err := LinkStat(scanner.Text())
		if err != nil {
			//report error, go no further...
			fmt.Fprintln(os.Stderr, "[x]", errors.Wrap(err, "LinkStat failed"))			
		} else {
			//TODO: correct scheme, e.g. for www. no http
			//TODO: if ftp, find alternative way to handle...
			//TODO: if response positive, populate LS further
			fmt.Fprintln(os.Stdout, "[x]", ls.ResponseCode, ls.ResponseText)
      }
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "[x]", errors.Wrap(err, "read failed"))
	}

}

func main() {
	looper()
}

