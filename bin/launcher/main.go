package main

import (
	"fmt"
	"go-writer/service"
	"time"

	logger "github.com/joaosoft/go-log/service"
)

func main() {
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
		log.Info(fmt.Sprintf("hello number %d\n", i))
	}
	fmt.Printf("sent!")

	// wait one minute to process...
	<-time.After(time.Minute * 1)
	quit <- true
}
