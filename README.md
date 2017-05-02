<div>
<p align="center">
<img id="logo" src="https://github.com/httpreserve/httpreserve/raw/master/src/images/httpreserve-logo.png" alt="httpreserve"/>
</p>
</div>

# httpreserve
[![Build Status](https://travis-ci.org/httpreserve/httpreserve.svg?branch=master)](https://travis-ci.org/httpreserve/httpreserve)
[![GoDoc](https://godoc.org/github.com/httpreserve/httpreserve?status.svg)](https://godoc.org/github.com/httpreserve/httpreserve)
[![Go Report Card](https://goreportcard.com/badge/github.com/httpreserve/httpreserve)](https://goreportcard.com/report/github.com/httpreserve/httpreserve)

A tool to check the status of a weblink and also see whether it is archived
in the [Internet Archive](https://archive.org/). 

Try it out here [httpreserve.info](http://httpreserve.info)

## Default Server

The library comes with a default server mode that can be configured for
POST and GET requests. POST by default. Default port is :2040 but this can
also be selected at runtime.

<img id="logo" src="https://github.com/httpreserve/httpreserve/raw/master/src/images/defaultserver.png" alt="httpreserve"/>

The default server can also be stood up as a web service. The API is
documented below. 

## Client

The httpreserve client is a separate application offering a broader range of
access methods. See: https://github.com/httpreserve/httpreserve-app

The client application is a work in progress. Stay tuned for more
information about its capabilities. 

## API

Primary entry point when the server is running:

*http://{httpreserve-ip-address}:{port}/httpreserve*

or 

*http://{httpreserve-ip-address}:{port}/save*

**GET** example:

* Return JSON struct with information about the service you requested:

    [http://<i>{httpreserve-ip-address}:{port}</i>/httpreserve?url=http://www.google.com&filename=filename.txt](http://httpreserve.info/httpreserve?url=http://www.google.com&filename=filename.txt)

* Manage a save request to the internet archive and return HTTPreserve struct:

    [http://<i>{httpreserve-ip-address}:{port}</i>/save?url=http://www.google.com&filename=filename.txt](http://httpreserve.info/httpreserve?url=http://www.google.com&filename=filename.txt)

**POST** example:

    Same access point, but encode url and filename in a <i>application/x-www-form-urlencoded</i> form.

**INFO** example: 

    TODO: Add some information to INFO response

**RETURN** value:

'application/json' struct to work with, e.g. 

      {
         "FileName": "",
         "AnalysisVersionNumber": "0.0.0",
         "AnalysisVersionText": "exponentialDK-httpreserve/0.0.0",
         "Link": "http://www.bbc.co.uk/",
         "ResponseCode": 200,
         "ResponseText": "OK",
         "ScreenShot": "{base64 encoded data}",
         "InternetArchiveLinkLatest": "http://web.archive.org/web/20170326191259/http://www0.bbc.co.uk/",
         "InternetArchiveLinkEarliest": "http://web.archive.org/web/19961221203254/http://www0.bbc.co.uk/",
         "InternetArchiveSaveLink": "http://web.archive.org/save/http://www.bbc.co.uk/",
         "InternetArchiveResponseCode": 200,
         "InternetArchiveResponseText": "OK",
         "Archived": true,
         "Error": false,
         "ErrorMessage": "",
         "StatsCreationTime": "4.506728277s"
      }

## Archiving Weblinks

* [Find and Connect Project:](http://www.findandconnectwrblog.info/2016/11/broken-links-broken-trust/) Nicola Laurent on the impact of broken links.
* [Binary Trees? Automatically Identifying the links between born digital records:](https://www.youtube.com/watch?v=Ked9GRmKlRw) I write about hyperlinks as a public record in own right when submitted as part of a documentary heritage.

### License

GNU General Public License Version 3. [Full Text](LICENSE)
