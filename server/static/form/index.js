// Simple AJAX function to make our demo page nice and clean
function httpreserve() {

	var key=""
	var value=""

	var elements = document.getElementById("abc").elements;
	var obj ={};
	for(var i = 0 ; i < elements.length ; i++){
		var item = elements.item(i);
		if (item.name == "url") {
			key = item.name;
			value = item.value;
		}
	}

	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (xmlhttp.readyState == XMLHttpRequest.DONE ) {
			if (xmlhttp.status == 200) {
				document.getElementById("httpreserve-analysis").innerHTML = xmlhttp.responseText;
			}
			else if (xmlhttp.status == 400) {
				alert('There was an error 400');
			}
			else {
				alert('something else other than 200 was returned');
			}
		}
	};

	//there?name=ferret
	//xmlhttp.open("POST", "httpreserve", true);
	//xmlhttp.setRequestHeader("Content-type","application/x-www-form-urlencoded");
	xmlhttp.open("GET", "httpreserve?url=http://example.com", true);
	xmlhttp.send(key + "=" + value);
}
