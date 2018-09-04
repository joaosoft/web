# webserver
[![Build Status](https://travis-ci.org/joaosoft/webserver.svg?branch=master)](https://travis-ci.org/joaosoft/webserver) | [![codecov](https://codecov.io/gh/joaosoft/webserver/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/webserver) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/webserver)](https://goreportcard.com/report/github.com/joaosoft/webserver) | [![GoDoc](https://godoc.org/github.com/joaosoft/webserver?status.svg)](https://godoc.org/github.com/joaosoft/webserver)

A simple web server. [UNDER DEVELOPMENT]

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Methods (HEAD, GET, POST, PUT, CONNECT, PATCH, DELETE, OPTIONS, TRACE)

>### Go
```
go get github.com/joaosoft/webserver
```

## Usage 
```
w, err := webserver.NewWebServer()
if err != nil {
    panic(err)
}

if err := w.Start(); err != nil {
    panic(err)
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
