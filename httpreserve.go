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

func testConnection (requrl string) (LinkStats, error) {
	var ls LinkStats
	var err error

	req, err := url.Parse(requrl)
	if err != nil {
		return ls, errors.Wrap(err, "url parse failed")
	}

	switch req.Scheme {
	case "ftp":
		/*_, err := handleftp(request)
		if err != nil {
			panic(err)
		}*/
	case "http":
		fallthrough
	case "https":
		ls, err = httpFromSimpleRequest(DefaultSimpleRequest(req))
		if err != nil {
			return ls, errors.Wrap(err, "handlehttp() failed")
		}
		return ls, nil
	case "":
		ls.Link = req
		return ls, errors.New(ERR_BLANK_PROTOCOL)		
	default:
		ls.Link = req
		return ls, errors.Wrap(errors.New(ERR_UNKNOWN_PROTOCOL), req.Scheme)
	}
	ls.Link = req
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
		fmt.Fprintln(os.Stderr, "[stdio]", errors.Wrap(err, "read failed"))
	}
}

func looper() {

	file, err := os.Open("link-examples/linkswork.txt") 
	if err != nil {
		fmt.Fprintln(os.Stderr, "[File Open]", errors.Wrap(err, "file open failed"))
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			//send through to our function tog get stats...
			link := scanner.Text()
			ls, err := LinkStat(link)
			if err != nil {
				//rreport error in some way...
				ls.ProtocolError = true
				switch err.Error() {
				case ERR_BLANK_PROTOCOL:
					ls.ProtocolErrorMessage = ERR_BLANK_PROTOCOL 
				case ERR_UNKNOWN_PROTOCOL:
					ls.ProtocolErrorMessage = ERR_UNKNOWN_PROTOCOL
				default:
					//TODO: Consider a log file
					//TODO: correct scheme, e.g. for www. add http://
					fmt.Fprintln(os.Stderr, "[LinkStat Fail]", errors.Wrap(err, "LinkStat failed"))
				}
			}
			//TODO: Positive or negative... populate LS
			ls, err = makeLinkStats(ls)
			if err != nil {
				fmt.Fprintln(os.Stderr, "[Make LinkStat]", err)
			}
			fmt.Fprintln(os.Stderr, "[LinkStat Returned]", ls.ResponseCode, ls.ResponseText, ls.InternetArchiveResponseCode, ls.InternetArchiveResponseText)
   	}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "[Scan Error]", errors.Wrap(err, "read http links failed"))
	}
}

func main() {

	looper()

	//server here
	//consider two servers
	//one - basic return LinkStats
	//two - assemble reports plus other information

	//oneoff save examples...

	//https://web.archive.org/save/http://www.bbc.com/news
	//https://web.archive.org/save/https://www.theguardian.com/international
	/*ls, err := handlehttp("https://web.archive.org/save/https://www.theguardian.com/international")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ls.header, ls.ResponseText, ls.ResponseCode)*/
}

