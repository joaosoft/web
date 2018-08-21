# builder
[![Build Status](https://travis-ci.org/joaosoft/builder.svg?branch=master)](https://travis-ci.org/joaosoft/builder) | [![codecov](https://codecov.io/gh/joaosoft/builder/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/builder) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/builder)](https://goreportcard.com/report/github.com/joaosoft/builder) | [![GoDoc](https://godoc.org/github.com/joaosoft/builder?status.svg)](https://godoc.org/github.com/joaosoft/builder)

A simple golang build and restart tool when some file of the project changes

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Rebuild
* Restart

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
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	w, err := builder.NewBuilder(builder.WithReloadTime(1))
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	if err := w.Start(&wg); err != nil {
		panic(err)
	}

	<-termChan
	wg.Add(1)
	if err := w.Stop(&wg); err != nil {
		panic(err)
	}
}
```


> Configuration file
```
{
  "builder": {
    "source": "main/main.go",
    "destination": "bin/builder",
    "reload_time": 1,
    "log": {
      "level": "error"
    }
  },
  "watcher": {
    "reload_time": 1,
    "dirs": {
      "watch":[ "." ],
      "excluded":[ "vendor", "bin" ],
      "extensions": [ "go", "json", "yml" ]
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
