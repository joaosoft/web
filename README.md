# watcher
[![Build Status](https://travis-ci.org/joaosoft/watcher.svg?branch=master)](https://travis-ci.org/joaosoft/watcher) | [![codecov](https://codecov.io/gh/joaosoft/watcher/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/watcher) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/watcher)](https://goreportcard.com/report/github.com/joaosoft/watcher) | [![GoDoc](https://godoc.org/github.com/joaosoft/watcher?status.svg)](https://godoc.org/github.com/joaosoft/watcher)

A simple cross-platform file watcher

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Multi directories
* Exclusions
* Extensions

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/watcher
```

## Usage 
This examples are available in the project at [watcher/main/main.go](https://github.com/joaosoft/watcher/tree/master/main/main.go)
```
import (
	github.com/joaosoft/watcher
	"fmt"
)

func main() {
	c := make(chan *service.Event)
	w, err := service.NewWatcher(service.WithEventChannel(c))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case event := <-c:
				fmt.Printf("received event %+v\n", event)
			}
		}
	}()

	if err := w.Start(); err != nil {
		panic(err)
	}
}
```


> Configuration file
```
{
  "watcher": {
    "host": "localhost:8001",
    "dirs": {
      "watch":[ "examples/" ],
      "excluded":[ "examples/test_2" ],
      "extensions": [ "go" ]
    },
    "log": {
      "level": "error"
    }
  },
  "manager": {
    "log": {
      "level": "error"
    }
  }
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
