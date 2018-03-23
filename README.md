# go-mapper
[![Build Status](https://travis-ci.org/joaosoft/go-mapper.svg?branch=master)](https://travis-ci.org/joaosoft/go-mapper) | [![Code Climate](https://codeclimate.com/github/joaosoft/go-mapper/badges/coverage.svg)](https://codeclimate.com/github/joaosoft/go-mapper)

Translates any struct to other data types.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## Convertions
* to map with key = path and value = the value

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/go-mapper/service
```

## Usage 
This examples are available in the project at [go-mapper/bin/launcher/main.go](https://github.com/joaosoft/go-mapper/tree/master/bin/launcher/main.go)

```go
type first struct {
	One   string
	Two   int
	Three map[string]string
	Four  Four
	Seven []string
	Eight []Four
}

type Four struct {
	Five string
	Six  int
}

type second struct {
	Eight []Four
}

obj1 := first{
    One:   "one",
    Two:   2,
    Three: map[string]string{"a": "1", "b": "2"},
    Four: Four{
        Five: "five",
        Six:  6,
    },
    Seven: []string{"a", "b", "c"},
    Eight: []Four{Four{Five: "5", Six: 66}},
}

log.Info("translate...")
mapper := gomapper.NewMapper(gomapper.WithLogger(log))
if translated, err := mapper.ToMap(obj1); err != nil {
    log.Error("error on translation!")
} else {
    log.Info("translated with success!")

    for key, value := range translated {
        fmt.Printf("%s: %+v\n", key, value)
    }
}

obj2 := second{
    Eight: []Four{Four{Five: "5", Six: 66}},
}
log.Info("translate...")
if translated, err := mapper.ToMap(obj2); err != nil {
    log.Error("error on translation!")
} else {
    log.Info("translated with success!")

    for key, value := range translated {
        fmt.Printf("%s: %+v\n", key, value)
    }
}
```

#### Result:
```javascript
{"prefixes":{"level":"info","service":"go-mapper","time":"2018-03-23 19:38:50:18"},"tags":{},"message":"translate...","fields":{}}
{"prefixes":{"level":"info","service":"go-mapper","time":"2018-03-23 19:38:50:18"},"tags":{},"message":"translated with success!","fields":{}}
three.{b}: 2
four.five: five
four.six: 6
seven.[0]: a
seven.[2]: c
eight.[0].six: 66
two: 2
three.{a}: 1
seven.[1]: b
eight.[0].five: 5
one: one
{"prefixes":{"level":"info","service":"go-mapper","time":"2018-03-23 19:38:50:18"},"tags":{},"message":"translate...","fields":{}}
{"prefixes":{"level":"info","service":"go-mapper","time":"2018-03-23 19:38:50:18"},"tags":{},"message":"translated with success!","fields":{}}
eight.[0].five: 5
eight.[0].six: 66
```

## Known issues


## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
