package httpreserve

import (
	"log"
	"net/http"
	"strings"
)

const templateFormMethod = "{{template.form.method}}"
const replacePOST = " method=\"post\" "
const replaceGET = " method=\"get\" "

// GetDefaultServerPage will print an example server page
// with the template filled in using the given method. Useful
// for folk who want to learn about this app. A bit of Ajax
// and a pretty decent way of encoding Base64 data hyperlinks.
func GetDefaultServerPage(method string) string {
	// A good default is to use POST as it's more secure
	// let users change to GET if they so wish,
	switch strings.ToUpper(method) {
	case http.MethodGet:
		log.Println("Default server HTTP method set to:", replaceGET)
		httpreservePages = strings.Replace(httpreservePages, templateFormMethod, replaceGET, 1)
		return httpreservePages
	case http.MethodPost:
		log.Println("Default server HTTP method set to:", replacePOST)
		httpreservePages = strings.Replace(httpreservePages, templateFormMethod, replacePOST, 1)
		return httpreservePages
	default:
		log.Printf("%s is not a valid HTTP method, setting to default: POST\n", method)
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
						if (xmlhttp.status == 200) {
							document.getElementById("httpreserve-analysis").innerHTML = xmlhttp.responseText;
						}
						else if (xmlhttp.status == 400) {
						   document.getElementById("httpreserve-analysis").innerHTML = '[WARNING] There was an error 400';
						}
						else {
						   document.getElementById("httpreserve-analysis").innerHTML = '[WARNING] something else other than 200 was returned';
						}
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
			figcaption { font-family: times new roman, arial, verdana; font-size: 22px; font-weight: bold; margin-bottom: 8px; }

			div.wrap { margin: 0 auto ; width:715px; }
			div.layout { margin-top: 50px }

			h4 { font-size: 16px;}

			input.link { display: block; margin: auto; width: 500px; }
			input.button  { display: block; margin: auto; }

			footer { display: block; bottom: 0; margin-bottom: 20px; position: fixed; }
		</style>
	</head>
	<body>
	<div class="wrap">
		<div class="layout">
		<p>
		<center>
			<figure>
			  	<figcaption>httpreserve</figcaption>
				<img src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmc
				vMjAwMC9zdmciIHdpZHRoPSI4IiBoZWlnaHQ9IjgiIHZpZXdCb3g9IjAgMCA4IDgiPg0KIC
				A8cGF0aCBkPSJNMCAwdjFoOHYtMWgtOHptNCAybC0zIDNoMnYzaDJ2LTNoMmwtMy0zeiIgLz4NCjwvc3ZnPg==" 
				width="60px" height="60px" alt="httpreserve"/>
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
		<footer>
		   A project by <a href="https://twitter.com/beet_keeper" alt="@beet_keeper on Twitter">@beet_keeper</a>
		   <br/>
		   On GitHub: <a href="https://github.com/exponential-decay/httpreserve" alt="httpreserve on GitHub">httpreserve</a>
		</footer>
		</div>
	</div>
	</body>
	</html> 

`
