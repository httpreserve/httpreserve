package main

import (
"net/http/httputil"
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

//hackable uri for saving pages
//https://web.archive.org/save/https://www.theguardian.com/politics/2017/mar/13/ian-mcewan-clarifies-remarks-likening-brexit-vote-third-reich

const USE_PROXY = false

const CONN_OKAY int8 = 0
const CONN_BAD int8 = 1

const GET = http.MethodGet
const HEAD = http.MethodHead

const USERAGENT = "exponentialDK-httpreserve/0.0.0"
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

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}

//proxy help:
//https://jannewmarch.gitbooks.io/network-programming-with-go-golang-/content/http/proxy_handling.html
//another example:
//http://stackoverflow.com/questions/40817784/access-https-via-http-proxy-with-basic-authentication
func handlehttp(request string, proxflag bool) (LinkStats, error) {

	var ls LinkStats
	var client = &http.Client{}

	linkurl, err := url.Parse(request)
	if err != nil {
		return ls, errors.Wrap(err, "parse request url failed")
	}

	req, err := http.NewRequest(HEAD, linkurl.String(), nil)
	if err != nil {
		return ls, errors.Wrap(err, "request generation failed")
	}
   req.Header.Add("User-Agent", USERAGENT)
   req.Header.Add("Range", BYTERANGE) 

	if proxflag {
		client, err = returnProxyClient(req)	
		if err != nil {
			return ls, errors.Wrap(err, "proxy header creation failed")
		}
	} 

	dump, _ := httputil.DumpRequest(req, false)
	fmt.Println("Request header:")
	fmt.Fprintln(os.Stdout, string(dump))

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

