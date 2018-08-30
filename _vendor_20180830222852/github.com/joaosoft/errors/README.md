# errors
[![Build Status](https://travis-ci.org/joaosoft/errors.svg?branch=master)](https://travis-ci.org/joaosoft/errors) | [![codecov](https://codecov.io/gh/joaosoft/errors/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/errors) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/errors)](https://goreportcard.com/report/github.com/joaosoft/errors) | [![GoDoc](https://godoc.org/github.com/joaosoft/errors?status.svg)](https://godoc.org/github.com/joaosoft/errors)

Error manager with error and caused by structure.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/errors
```

## Usage 
This examples are available in the project at [example_test.gp](https://github.com/joaosoft/errors/tree/master/example_test.go)
```go
func TestExampleSimple(t *testing.T) {

	err := New("1", "erro 1")
	err.Add(New("2", "erro 2"))
	err.Add(New("3", "erro 3"))

	fmt.Printf("Error: %s, Cause: %s", err.String(), err.Cause())

	assert.Equal(t, err.String(), `{"previous":{"previous":{"code":"1","error":"erro 1"},"code":"2","error":"erro 2"},"code":"3","error":"erro 3"}`)
	assert.Equal(t, err.Cause(), `'erro 3', caused by 'erro 2', caused by 'erro 1'`)
}

func TestExampleList(t *testing.T) {

	var errs ListErr
	errs.Add(New("1", "erro 1"))
	errs.Add(New("2", "erro 2"))
	errs.Add(New("3", "erro 3"))

	fmt.Printf("Errors: %s", errs.String())

	assert.Equal(t, errs.String(), `[{"code":"1","error":"erro 1"},{"code":"2","error":"erro 2"},{"code":"3","error":"erro 3"}]`)
}
```

## Known issues


## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
