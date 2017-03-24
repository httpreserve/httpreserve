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