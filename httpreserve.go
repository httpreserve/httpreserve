package main

import (
	"os"
	"fmt"
	"bufio"
	"net/url"
	"github.com/pkg/errors"
   )

//ftp: ftp://exponentialdecay.co.uk/
//http; http://exponentialdecay.co.uk
//https: https://github.com/exponential-decay

//hackable uri for archive.org, finds closest to date set, pre or post
//internet archive uri: http://web.archive.org/web/20161104020243/http://exponentialdecay.co.uk/#

//hackable uri for saving pages
//https://web.archive.org/save/https://www.theguardian.com/politics/2017/mar/13/ian-mcewan-clarifies-remarks-likening-brexit-vote-third-reich

const USE_PROXY = false

const CONN_OKAY int8 = 0
const CONN_BAD int8 = 1

const USERAGENT = "exponentialDK-httpreserve/0.0.0"
const BYTERANGE = "bytes=0-0"

func testConnection (request string) (LinkStats, error) {
	var ls LinkStats
	u, err := url.Parse(request)
	if err != nil {
		return ls, errors.Wrap(err, "url parse failed")
	}

	switch u.Scheme {
	case "ftp":
		/*_, err := handleftp(request)
		if err != nil {
			panic(err)
		}*/
	case "http":
		fallthrough
	case "https":
		ls, err := handlehttp(request, USE_PROXY)
		if err != nil {
			return ls, errors.Wrap(err, "handlehttp() failed")
		}
		return ls, nil
	}
	return ls, nil
}

func LinkStat(url string) (LinkStats, error) {
	var ls LinkStats
	ls, err := testConnection(url)
	return ls, err
}

//might be useful... stdin
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

	file, err := os.Open("link-examples/linkswork.txt") 
	if err != nil {
		fmt.Fprintln(os.Stderr, "[x]", errors.Wrap(err, "file open failed"))
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			//send through to our function tog get stats...
			ls, err := LinkStat(scanner.Text())
			if err != nil {
				//report error, go no further...
				fmt.Fprintln(os.Stderr, "[x]", errors.Wrap(err, "LinkStat failed"))			
			} else {
				//TODO: correct scheme, e.g. for www. no http
				//TODO: if ftp, find alternative way to handle...
				//TODO: if response positive, populate LS further
				fmt.Fprintln(os.Stdout, "[success]", ls.ResponseCode, ls.ResponseText)
	      }
   	}
   	fmt.Println("----\n\n")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "[x]", errors.Wrap(err, "read failed"))
	}
}

func main() {

	looper()

	//oneoff save examples...

	//https://web.archive.org/save/http://www.bbc.com/news
	//https://web.archive.org/save/https://www.theguardian.com/international
	/*ls, err := handlehttp("https://web.archive.org/save/https://www.theguardian.com/international")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ls.header, ls.ResponseText, ls.ResponseCode)*/
}

