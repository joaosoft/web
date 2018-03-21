# go-writer
[![Build Status](https://travis-ci.org/joaosoft/go-writer.svg?branch=master)](https://travis-ci.org/joaosoft/go-writer) | [![Code Climate](https://codeclimate.com/github/joaosoft/go-writer/badges/coverage.svg)](https://codeclimate.com/github/joaosoft/go-writer)

A starting project with writer interface implementations.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* file writer

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/go-writer/service
```

## Interface 
```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

## Usage 
This examples are available in the project at [go-writer/bin/launcher/main.go](https://go-writer/tree/master/bin/launcher/main.go)

```go
quit := make(chan bool)
//
// file writer
writer := gowriter.NewFileWriter(
    gowriter.WithDirectory("./testing"),
    gowriter.WithFileName("dummy_"),
    gowriter.WithFileMaxMegaByteSize(1),
    gowriter.WithFlushTime(time.Second),
    gowriter.WithQuitChannel(quit),
)

// logger
log := logger.NewLog(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormatHandler(logger.JsonFormatHandler),
    logger.WithWriter(writer)).WithPrefixes(map[string]interface{}{
    "level":   logger.LEVEL,
    "time":    logger.TIME,
    "service": "go-writer"})

fmt.Printf("send...")
for i := 1; i < 100000; i++ {
    log.Info(fmt.Sprintf("ola %d\n", i))
}
fmt.Printf("sent!")

// wait one minute to process...
<-time.After(time.Minute * 1)
```

## Known issues


## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
