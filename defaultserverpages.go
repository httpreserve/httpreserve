package httpreserve

import (
	"net/http"
	"strings"
)

const templateFormMethod = "{{template.form.method}}"
const replacePOST = " method=\"post\" "
const replaceGET = " method=\"get\" "

const loadTMP = "{{ LOADING }}"

// GetDefaultServerPage will print an example server page
// with the template filled in using the given method. Useful
// for folk who want to learn about this app. A bit of Ajax
// and a pretty decent way of encoding Base64 data hyperlinks.
func GetDefaultServerPage(method string) string {

	// Make loading gif html page friendly...
	loader := strings.Replace(loading, "\n", "", -1)
	loader = strings.Trim(loader, " ")

	// First update the load image
	httpreservePages = strings.Replace(httpreservePages, loadTMP, loader, 1)

	// A good default is to use POST as it's more secure
	// let users change to GET if they so wish,
	switch strings.ToUpper(method) {
	case http.MethodGet:
		httpreservePages = strings.Replace(httpreservePages, templateFormMethod, replaceGET, 1)
		return httpreservePages
	case http.MethodPost:
		httpreservePages = strings.Replace(httpreservePages, templateFormMethod, replacePOST, 1)
		return httpreservePages
	default:
		httpreservePages = strings.Replace(httpreservePages, templateFormMethod, replacePOST, 1)
		return httpreservePages
	}
}

// Web entry point for our default server for demo purposes.
var httpreservePages = `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="utf-8">
        <title>httpreserve</title>

        <meta name="description" content="Retrieve Internet Archive links and metadata easily with HTTPreserve">
        <meta name="og:title" property="og:title" content="HTTPreserve - HTTPpreservation of web pages">

        <link rel="canonical" href="https://httpreserve.info">

        <link id="favicon" rel="icon" type="image/x-icon"
        href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAABm
        JLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAA
        B3RJTUUH4QMZAiAVgKAUZQAAAGlJREFUWMPt1jkSgDAMQ9GYk/
        vmoqIM4C1qpAPkv5k0NgBYxF2LPAE0rTQADsCZ8WfOjJ9FbOJn
        EB/xWcTP+AwiGO9FJOM9iGK8hmiK5xDN8RhiKL5FWOJrXm9IMw
        u9qZNMAAEEEIAOuAEmTWnhcv1r2AAAAABJRU5ErkJggg=="/>

        <script type="text/javascript">
            // Simple AJAX function to make our demo page nice and clean
            function httpreserve() {

                var myform = document.getElementById("httpreserve-form");

                var key="";
                var value="";

                var method = myform.method;
                var elements = myform.elements;

                var obj ={};
                for(var i = 0 ; i < elements.length ; i++){
                    var item = elements.item(i);
                    if (item.name == "url") {
                        key = item.name;
                        value = item.value;
                    }
                }

                // Add HTTP:// by default if not in string
                value = addhttp(value)

                var xmlhttp = new XMLHttpRequest();
                xmlhttp.onreadystatechange = function() {
                    if (xmlhttp.readyState == XMLHttpRequest.DONE ) {

                        document.getElementById("loader").innerHTML = "";

                        if (xmlhttp.status == 200) {
                            document.getElementById("httpreserve-analysis").innerHTML = xmlhttp.responseText;
                        }
                        else if (xmlhttp.status == 400) {
                           document.getElementById("httpreserve-analysis").innerHTML = '[WARNING] There was an error 400';
                        }
                        else {
                           document.getElementById("httpreserve-analysis").innerHTML = '[WARNING] something else other than 200 was returned';
                        }
                    } else {
                        document.getElementById("httpreserve-analysis").innerHTML = "";
                        document.getElementById("loader").innerHTML = {{ LOADING }};
                    }
                };

                if (method.toLowerCase() == "post") {
                    xmlhttp.open("POST", "httpreserve", true);
                    xmlhttp.setRequestHeader("Content-type","application/x-www-form-urlencoded");
                    xmlhttp.send(key + "=" + value);
                    return;
                }

                if (method.toLowerCase() == "get") {
                    xmlhttp.open("GET", "httpreserve?" + key + "=" + value, true);
                    xmlhttp.send(key + "=" + value);
                    return;
                }

                document.getElementById("httpreserve-analysis").innerHTML = '[WARNING] issue parsing the form in JavaScript';
            }

            // Search for HTTP:// in our string for usability of tool
            // will fail if we haven't a protocol.
            function addhttp(value) {
                if (value.indexOf("http://") != -1) {
                    return value
                }
                if (value.indexOf("https://") != -1) {
                    return value
                }
                if (value.indexOf("ftp://") != -1) {
                    return value
                }
                var newurl = "http://".concat(value)
                return newurl
            }

        </script>
        <style>
            body { min-width:750px; font-family: arial, verdana; font-size: 14px; margin-bottom: 20px}
            figure { font-family: arial, helvetica, verdana; font-size: 40px; font-weight: bold; margin-bottom: 8px; }
            h1 { font-size: 40px; }

            figure.loading { font-family: arial, verdana; font-size: 8px; margin-top: -2px; font-weight: normal; }

            div.wrap { margin: 0 auto ; width:715px; }
            div.layout { margin-top: 50px; min-height: 100%; height: 250px; }

            h4 { font-size: 16px; margin-bottom: -2px;}

            input.link { display: block; margin: auto; width: 500px; }
            input.button  { display: block; margin: auto; width: 200px; height: 30px; font-size: 20px; font-family: courier; font-weight: bold; }

            pre.analysis { font-size: 14px; }

         /*use push to position footer more usefully on screen if necessary*/
         div.push { height: 540px; min-height: 540px; }

         div.footer { height: 50px; margin: 0 auto ; width:200px; text-align: center; }

            /* Rotate via: https://linuxhint.com/rotate-animation-css/ */
            .rotate {
              animation: rotation 5s infinite linear;
            }

            img.rotate { margin-top: 50px; }

            @keyframes rotation {
              from {
                transform: rotateY(0deg);
              }
              to {
                transform: rotateY(359deg);
              }
            }

        </style>

    </head>
    <body>
    <div class="wrap">
        <div class="layout">
        <p>
        <center>
            <figure>
                <figcaption style="font-family: helvetica; arial, verdana; margin-bottom: 8px"><h1>httpreserve</h1></figcaption>
                <img src=" ` + httpreserveImage + `"
                width="80px" height="80px" alt="httpreserve"/>
            </figure>
        </center>
        </p>
        <h4>Enter a URL:</h4>
        <p>
        This service will tell you if an internet archive link exists and give you a mechanism to save it as you wish.
        <br/>
        <br/>
        <form action="javascript:httpreserve()" ` + templateFormMethod + ` id="httpreserve-form">
           <input class="link" type="text" name="url" value="http://example.com/">
           <br/>
           <input class="button" type="submit" value="httpreserve">
        </form>
        <br/><br/>
        <pre class="analysis" id="httpreserve-analysis"></div>
        </p>
        <div>
        <center>
        <div id="loader"></div>
        </center>
        </div>
        <div class="push">&nbsp;</div>
        </div>
      <div class="footer" id="footer">
        A project by <a href="https://twitter.com/beet_keeper" alt="@beet_keeper on Twitter">@beet_keeper</a>
        <br/>
        On GitHub: <a href="https://github.com/exponential-decay/httpreserve" alt="httpreserve on GitHub">httpreserve</a>
    </div>
    </div>
    </body>
    </html>

`

var loading string = "'<figure><img class=\"rotate\" src=\"" + httpreserveImage + "\" " +
	"width=\"100px\" height=\"100px\" alt=\"loading\"/>" +
	"<figcaption class=\"loading\" style=\"font-family: helvetica; arial, verdana;\">" +
	"processing</figcaption></figure>'"
