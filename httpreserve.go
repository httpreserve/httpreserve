package main

import (
		"os"
      "net/http"
      "fmt"
		//"io/ioutil"
      //"bufio"
   )

//ftp: ftp://exponentialdecay.co.uk/
//http; http://exponentialdecay.co.uk
//https: https://github.com/exponential-decay

//hackable uri for archive.org, finds closest to date set, pre or post
//internet archive uri: http://web.archive.org/web/20161104020243/http://exponentialdecay.co.uk/#

const CONN_OKAY int8 = 0
const CONN_BAD int8 = 1

const GET = http.MethodGet

func testConnection (request string) int8 {

   conn := CONN_OKAY
	req, err := http.NewRequest(GET, request, nil) 
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: error creating request,", err)
      os.Exit(1)
	}

   req.Header.Add("User-Agent", "@exponentialDK Digital Preservation of HTTP in documentary heritage.")
   req.Header.Add("Range", "bytes=0-0") 

	fmt.Println(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
      conn = CONN_BAD
		return conn
	}

	fmt.Println(resp.Status)
	fmt.Println(http.StatusText(resp.StatusCode))
	fmt.Println("")
	fmt.Println(resp.Header)
	fmt.Println("")

	//body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// display content to screen ... save this to a HTML file and view the file with browser ;-)
	//fmt.Println(string(body))

   return conn
}

func main() {
	//testing accept-range...
	a := testConnection("https://raw.githubusercontent.com/exponential-decay/the-format-registry/master/LICENSE")
	fmt.Println(a)
	fmt.Println("xxx")

	//404 returned if not there... (useful)
	a = testConnection("http://web.archive.org/web/20161104020243/http://exponentialdecayxxxx.co.uk/#")
	fmt.Println(a)
	fmt.Println("xxx")


}

