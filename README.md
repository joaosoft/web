# go-mapper
[![Build Status](https://travis-ci.org/joaosoft/go-mapper.svg?branch=master)](https://travis-ci.org/joaosoft/go-mapper) | [![codecov](https://codecov.io/gh/joaosoft/go-mapper/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/go-mapper) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/go-mapper)](https://goreportcard.com/report/github.com/joaosoft/go-mapper) | [![GoDoc](https://godoc.org/github.com/joaosoft/go-mapper?status.svg)](https://godoc.org/github.com/joaosoft/go-mapper/service)

Translates any struct to other data types.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## Convertions
* to map with key = path and value = the value
* to string with path = value

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
type First struct {
	One   string            `json:"one"`
	Two   int               `json:"two"`
	Three map[string]string `json:"three"`
	Four  Four              `json:"four"`
	Seven []string          `json:"seven"`
	Eight []Four            `json:"eight"`
}

type Four struct {
	Five string `json:"five"`
	Six  int    `json:"six"`
}

type Second struct {
	Eight []Four          `json:"eight"`
	Nine  map[Four]Second `json:"nine"`
}

obj1 := First{
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
obj2 := Second{
    Eight: []Four{Four{Five: "5", Six: 66}},
    Nine:  map[Four]Second{Four{Five: "111", Six: 1}: Second{Eight: []Four{Four{Five: "222", Six: 2}}}},
}
```
#### Convert struct to string 
```go
fmt.Println(":::::::::::: STRUCT ONE")
mapper := gomapper.NewMapper(gomapper.WithLogger(log))
if translated, err := mapper.String(obj1); err != nil {
    log.Error("error on translation!")
} else {
    fmt.Println(translated)
}

fmt.Println(":::::::::::: STRUCT TWO")
if translated, err := mapper.String(obj2); err != nil {
    log.Error("error on translation!")
} else {
    fmt.Println(translated)
}
```

##### Result:
```javascript
:::::::::::: STRUCT ONE
One: one
Two: 2
Three
  {a}: 1
  {b}: 2
Four
  Five: five
  Six: 6
Seven
  [0]: a
  [1]: b
  [2]: c
Eight
  [0]
    Five: 5
    Six: 66

:::::::::::: STRUCT TWO
Eight
  [0]
    Five: 5
    Six: 66
Nine
  {Five: 111 Six: 1}
    Eight
      [0]
        Five: 222
        Six: 2
    Nine

:::::::::::: JSON STRING OF STRUCT ONE
{seven}
  [0]: a
  [1]: b
  [2]: c
{eight}
  [0]
    {five}: 5
    {six}: 66
{one}: one
{two}: 2
{three}
  {a}: 1
  {b}: 2
{four}
  {five}: five
  {six}: 6
```

#### Convert struct to map 
```go
fmt.Println(":::::::::::: STRUCT ONE")
mapper := gomapper.NewMapper(gomapper.WithLogger(log))
if translated, err := mapper.Map(obj1); err != nil {
    log.Error("error on translation!")
} else {
    for key, value := range translated {
        fmt.Printf("%s: %+v\n", key, value)
    }
}

fmt.Println(":::::::::::: STRUCT TWO")
if translated, err := mapper.Map(obj2); err != nil {
    log.Error("error on translation!")
} else {
    for key, value := range translated {
        fmt.Printf("%s: %+v\n", key, value)
    }
}
```

##### Result:
```javascript
:::::::::::: STRUCT ONE
Four.Six: 6
Seven.[0]: a
Seven.[1]: b
Seven.[2]: c
Eight.[0].Six: 66
One: one
Four.Five: five
Three.{b}: 2
Eight.[0].Five: 5
Two: 2
Three.{a}: 1

:::::::::::: STRUCT TWO
Eight.[0].Five: 5
Eight.[0].Six: 66
Nine.{Five=111,Six=1}.Eight.[0].Five: 222
Nine.{Five=111,Six=1}.Eight.[0].Six: 2

:::::::::::: JSON STRING OF STRUCT ONE
{four}.{five}: five
{seven}.[0]: a
{seven}.[1]: b
{two}: 2
{three}.{b}: 2
{one}: one
{three}.{a}: 1
{four}.{six}: 6
{seven}.[2]: c
{eight}.[0].{five}: 5
{eight}.[0].{six}: 66
```

## Known issues


## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
