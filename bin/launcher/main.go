package main

import (
	"encoding/json"
	"fmt"
	"go-mapper/service"
	"os"

	logger "github.com/joaosoft/go-log/service"
	writer "github.com/joaosoft/go-writer/service"
)

var log = logger.NewLog(
	logger.WithLevel(logger.InfoLevel),
	logger.WithFormatHandler(writer.JsonFormatHandler),
	logger.WithWriter(os.Stdout)).WithPrefixes(map[string]interface{}{
	"level":   logger.LEVEL,
	"time":    logger.TIME,
	"service": "go-mapper"})

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

func main() {
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

	fmt.Println(":::::::::::: STRUCT TO MAP ::::::::::::")

	fmt.Println("\n:::::::::::: STRUCT ONE")

	mapper := gomapper.NewMapper(gomapper.WithLogger(log))
	if translated, err := mapper.Map(obj1); err != nil {
		log.Error("error on translation!")
	} else {
		for key, value := range translated {
			fmt.Printf("%s: %+v\n", key, value)
		}
	}

	fmt.Println("\n:::::::::::: STRUCT TWO")

	if translated, err := mapper.Map(obj2); err != nil {
		log.Error("error on translation!")
	} else {
		for key, value := range translated {
			fmt.Printf("%s: %+v\n", key, value)
		}
	}

	fmt.Println("\n:::::::::::: JSON STRING OF STRUCT ONE")

	bytesObj1, _ := json.Marshal(obj1)
	var convObj1 interface{}
	json.Unmarshal(bytesObj1, &convObj1)
	if translated, err := mapper.Map(convObj1); err != nil {
		log.Error("error on translation!")
	} else {
		for key, value := range translated {
			fmt.Printf("%s: %+v\n", key, value)
		}
	}

	fmt.Println("\n\n:::::::::::: STRUCT TO STRING ::::::::::::")

	fmt.Println("\n:::::::::::: STRUCT ONE")
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

	fmt.Println(":::::::::::: JSON STRING OF STRUCT ONE")
	bytesObj1, _ = json.Marshal(obj1)
	json.Unmarshal(bytesObj1, &convObj1)
	fmt.Println("STEING:" + string(bytesObj1))
	if translated, err := mapper.String(convObj1); err != nil {
		log.Error("error on translation!")
	} else {
		fmt.Println(translated)
	}
}
