# go-log
[![Build Status](https://travis-ci.org/joaosoft/go-log.svg?branch=master)](https://travis-ci.org/joaosoft/go-log) | [![Code Climate](https://codeclimate.com/github/joaosoft/go-log/badges/coverage.svg)](https://codeclimate.com/github/joaosoft/go-log)

A simplified logger that allows you to add complexity depending of your requirements.
After a read of the project https://gitlab.com/vredens/go-logger extracted some concepts like allowing to add tags and fields to logger infrastructure. 

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* formatted messages
* prefixes
* tags
* fields

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/go-log/service
```

## Interface 
```go
type Log interface {
	SetLevel(level Level)

	With(prefixes, tags, fields map[string]interface{}) Log
	WithPrefixes(prefixes map[string]interface{}) Log
	WithTags(tags map[string]interface{}) Log
	WithFields(fields map[string]interface{}) Log

	Debug(message interface{})
	Info(message interface{})
	Warn(message interface{})
	Error(message interface{})

	Debugf(format string, arguments ...interface{})
	Infof(format string, arguments ...interface{})
	Warnf(format string, arguments ...interface{})
	Errorf(format string, arguments ...interface{})
}
```

## Usage 
This examples are available in the project at [go-log/bin/launcher/main.go](https://go-log/tree/master/bin/launcher/main.go)

```go
//
// log to text
fmt.Println(":: LOG TEXT")
log := golog.NewLog(
    golog.WithLevel(golog.InfoLevel), 
    golog.WithFormatHandler(golog.TextFormatHandler), 
    golog.WithWriter(os.Stdout)).
        With(
            map[string]interface{}{"level": golog.LEVEL, "time": golog.TIME}, 
            map[string]interface{}{"service": "log"}, 
            map[string]interface{}{"name": "joão"})

// logging...
log.Error("isto é uma mensagem de error")
log.Info("isto é uma mensagem de info")
log.Debug("isto é uma mensagem de debug")

fmt.Println("--------------")
<-time.After(time.Second)

//
// log to json
fmt.Println(":: LOG JSON")
log = golog.NewLog(
    golog.WithLevel(golog.InfoLevel),
    golog.WithFormatHandler(golog.JsonFormatHandler),
    golog.WithWriter(os.Stdout)).
    With(
    map[string]interface{}{"level": golog.LEVEL, "time": golog.TIME},
    map[string]interface{}{"service": "log"},
    map[string]interface{}{"name": "joão"})

// logging...
log.Errorf("isto é uma mensagem de error %s", "hello")
log.Infof("isto é uma  mensagem de info %s ", "hi")
log.Debugf("isto é uma mensagem de debug %s", "ehh")
```

###### Output 

```javascript
:: LOG TEXT
{prefixes:map[level:error time:2018-03-20 02:47:21] tags:map[service:log] message:isto é uma mensagem de error fields:map[name:joão]}
{prefixes:map[level:info time:2018-03-20 02:47:21] tags:map[service:log] message:isto é uma mensagem de info fields:map[name:joão]}
--------------
:: LOG JSON
{"prefixes":{"level":"error","time":"2018-03-20 02:47:22"},"tags":{"service":"log"},"message":"isto é uma mensagem de error hello","fields":{"name":"joão"}}
{"prefixes":{"level":"info","time":"2018-03-20 02:47:22"},"tags":{"service":"log"},"message":"isto é uma  mensagem de info hi ","fields":{"name":"joão"}}
```

## Known issues
* all the maps do not guarantee order of the items! 


## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
