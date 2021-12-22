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
            body { min-width:750px; font-family: arial, verdana; font-size: 10px; margin-bottom: 20px}
            figure { font-family: times new roman, arial, verdana; font-size: 20px; font-weight: bold; margin-bottom: 8px; }

            figure.loading { font-family: arial, verdana; font-size: 8px; margin-top: -2px; font-weight: normal; }

            div.wrap { margin: 0 auto ; width:715px; }
            div.layout { margin-top: 50px; min-height: 100%; height: 250px; }

            h4 { font-size: 16px; margin-bottom: -2px;}

            input.link { display: block; margin: auto; width: 500px; }
            input.button  { display: block; margin: auto; }

         /*use push to position footer more usefully on screen if necessary*/
         div.push { height: 340px; min-height: 340px; }

         div.footer { height: 50px; margin: 0 auto ; width:200px; text-align: center; }
        </style>
    </head>
    <body>
    <div class="wrap">
        <div class="layout">
        <p>
        <center>
            <figure>
                <figcaption style="font-family: helvetica; arial, verdana; margin-bottom: 8px">httpreserve</figcaption>
                <img src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmc
                vMjAwMC9zdmciIHdpZHRoPSI4IiBoZWlnaHQ9IjgiIHZpZXdCb3g9IjAgMCA4IDgiPg0KIC
                A8cGF0aCBkPSJNMCAwdjFoOHYtMWgtOHptNCAybC0zIDNoMnYzaDJ2LTNoMmwtMy0zeiIgLz4NCjwvc3ZnPg=="
                width="50px" height="50px" alt="httpreserve"/>
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
        <pre id="httpreserve-analysis"></div>
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

const loading = `
'<figure><img src="data:image/gif;base64,R0lGODlhKgAqAPcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEB
AQEBAQICAgQEBAUFBQcHBwkJCQwMDA8PDxISEhYWFhkZGRoaGhsbGxwcHB0dHR4eHh8fHyAgICEh
ISIiIiMjIyQkJCUlJSYmJicnJygoKCkpKSoqKisrKywsLC0tLS4uLi8vLzAwMDExMTIyMjMzMzQ0
NDU1NTY2Njc3Nzk5OTo6Ojs7Ozw8PD09PT09PT4+Pj8/P0BAQEBAQEFBQUJCQkJCQkNDQ0REREVF
RUZGRkhISEpKSktLS0xMTE1NTU5OTk9PT1BQUFFRUVJSUlNTU1RUVFVVVVZWVldXV1hYWFlZWVpa
WltbW1xcXF1dXV5eXl9fX2BgYGFhYWJiYmNjY2RkZGVlZWZmZmdnZ2hoaGlpaWpqamtra2xsbG1t
bW5ubm9vb3BwcHFxcXJycnNzc3R0dHV1dXZ2dnd3d3h4eHl5eXp6enx8fH9/f4KCgoSEhIaGhoeH
h4mJiYqKiouLi4yMjIyMjI2NjY2NjY2NjY2NjY6Ojo6Ojo6Ojo6Ojo6Ojo6Ojo6Ojo6Ojo+Pj4+P
j5CQkJCQkJGRkZGRkZKSkpOTk5SUlJWVlZWVlZaWlpeXl5mZmZqampycnJ2dnZ+fn6GhoaOjo6am
pqenp6ioqKmpqaqqqqurq6ysrK2tra6urq+vr7CwsLGxsbKysrOzs7S0tLW1tba2tre3t7i4uLm5
ubq6uru7u7y8vL29vb6+vr+/v8DAwMHBwcLCwsPDw8TExMXFxcbGxsfHx8jIyMnJycrKysvLy87O
ztHR0dTU1NfX19ra2tvb2/T09Pr6+vz8/P39/f7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+
/v7+/v//////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////yH/C05FVFNDQVBFMi4wAwEA
AAAh+QQJAwDwACwAAAAAKgAqAAAIlgDhCRxIsKDBgwgTKlzIsKHDhxAjSpxIsaLFixgTQnPEsaNH
aBkLojpAsqRJVCEJjjTJEmVKgStZlnT5sqbNmzgVTst5MBAfngZJAiX4jGSzoQKDkNyBFJ7JnUAr
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
mguSJ2wmVWaBZAViUwGZ5DMViMkbU1FRIknJFNZwJO8KzKuX712/WAFPFdyUcFLDRREH5ai3sePH
kCNLphgQACH5BAkDAOEALAAAAAAqACoAAAiWAMMJHEiwoMGDCBMqXMiwocOHECNKnEixosWLGBWC
UsSxI0dQGQ9uAECyJMkNIQ2ONFkSZUqCK1kCcPmyps2bOHNOHOVJp0EMFnwSlEYymlCBeEjCORru
AsmgQqOZfCZ0isklQmVa04lMJjCdzUaNIimW2VGSTAWiTbuWaduzANKGe5s1Llu7bvEe5Si3r9+/
gAMLbhgQADsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=" width="30px" height="30px" alt="loading"/>
<figcaption class="loading" style="font-family: helvetica; arial, verdana;">processing</figcaption></figure>'
`
