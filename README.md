# builder
[![Build Status](https://travis-ci.org/joaosoft/builder.svg?branch=master)](https://travis-ci.org/joaosoft/builder) | [![codecov](https://codecov.io/gh/joaosoft/builder/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/builder) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/builder)](https://goreportcard.com/report/github.com/joaosoft/builder) | [![GoDoc](https://godoc.org/github.com/joaosoft/builder?status.svg)](https://godoc.org/github.com/joaosoft/builder)

A simple golang rebuild when some file changes

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Rebuild

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/builder
```

## Usage 
This examples are available in the project at [builder/main/main.go](https://github.com/joaosoft/builder/tree/master/main/main.go)
```
import (
	github.com/joaosoft/builder
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
  "builder": {
    "log": {
      "level": "error"
    }
  },
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
