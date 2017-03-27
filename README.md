<div>
<p align="center">
<img id="logo" src="https://github.com/exponential-decay/httpreserve/raw/master/src/images/httpreserve-logo.png" alt="httpreserve"/>
</p>
</div>

# httpreserve
[![Build Status](https://travis-ci.org/exponential-decay/httpreserve.svg?branch=master)](https://travis-ci.org/exponential-decay/httpreserve)
[![GoDoc](https://godoc.org/github.com/exponential-decay/httpreserve?status.svg)](https://godoc.org/github.com/exponential-decay/httpreserve)
[![Go Report Card](https://goreportcard.com/badge/github.com/exponential-decay/httpreserve)](https://goreportcard.com/report/github.com/exponential-decay/httpreserve)

Placeholder text to describe tool further. 

# Default Server

The library comes with a default servere mode that can be configured for
POST and GET requests. POST by default. Default port is :2040 but this can
also be selected at runtime.

<img id="logo" src="https://github.com/exponential-decay/httpreserve/raw/master/src/images/defaultserver.png" alt="httpreserve"/>

The default server can also be stood up as a web service. The API is
documented below. 

# Client

The httpreserve client is a separate application offering a broader range of
access methods. See: https://github.com/exponential-decay/httpreserve-app

# API

Primary entry point when the server is running:

http://{httpreserve-ip-address}:{port}/httpreserve

GET example:

http://{httpreserve-ip-address}:{port}/httpreserve?url=http://www.google.com&filename=filename.txt

POST example:

Same access point, but encode url and filename in a <i>application/x-www-form-urlencoded</i> form.

RETURN valie:

### See Also

* [Find and Connect Project:](http://www.findandconnectwrblog.info/2016/11/broken-links-broken-trust/) Nicola Laurent on the impact of broken links.
* [Binary Trees? Automatically Identifying the links between born digital records.](https://www.youtube.com/watch?v=Ked9GRmKlRw) I write about hyperlinks as a public record in own right when submitted as part of a documentary heritage.

### License

GNU General Public License Version 3. [Full Text](LICENSE)
