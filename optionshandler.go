package httpreserve

import (
	"fmt"
	"net/http"
	"strings"
)

const welcome = `
_ _ ______   _________  _____   _______    
| | ||___|   |   |  ||\/||___    | |  |    
|_|_||___|___|___|__||  ||___    | |__|....
                                           
_  __________ _____________________  _____ 
|__| |  | |__]|__/|___[__ |___|__/|  ||___ 
|  | |  | |   |  \|______]|___|  \ \/ |___ 

`

var byline = "A service to help you audit the state of your weblinks."

var instructions = `
There are two things you can do with this service, these are documented below.

See the status of your weblink:

   curl -i -X GET http://httpreserve.info/httpreserve?url=http://www.example.com

Manage the transaction with the wayback machine to save your link:

   curl -i -X GET http://httpreserve.info/save?url=http://www.example.com

POST will also work if you encode your form:

   application/x-www-form-urlencoded

Filename is encoded in the URL to help folks audit their repositories.
It is expected they will extract the links contained deep in their documentary
heritage, word docs, pdf, wordperfect, etc. using a tool like tikalinkextractor

   https://github.com/httpreserve/tikalinkextract

And then run it through httpreserve to report on the links still standing since
the advent of the web and the intersection with the paper-paradigm for office
productivity.

httpreserve.info is a demo server running on the cheapest VPS option available
at Linode.com which means it's really just a demo for now but feel free to wrap
it into your scripts.

See workbench app for a more powerful theoretical way to use this tool with
bulk lists of links.

   https://github.com/httpreserve/workbench

A more practical implementation is at:

   https://github.com/httpreserve/linkstat

The full suite of tools can be found at httpreserve suite.

   https://github.com/httpreserve/
`

func handleOptions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n\n", strings.Trim(welcome, " \n"))
	fmt.Fprintf(w, "%s\n\n", byline)
	fmt.Fprintf(w, "%s\n\n", strings.Trim(instructions, " \n"))
}
