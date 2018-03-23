package main

import (
	"fmt"
	"go-mapper/service"
	"os"

	logger "github.com/joaosoft/go-log/service"
)

var log = logger.NewLog(
	logger.WithLevel(logger.InfoLevel),
	logger.WithFormatHandler(logger.JsonFormatHandler),
	logger.WithWriter(os.Stdout)).WithPrefixes(map[string]interface{}{
	"level":   logger.LEVEL,
	"time":    logger.TIME,
	"service": "go-mapper"})

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

func main() {
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
}
