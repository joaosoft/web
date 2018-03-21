package main

import (
	"fmt"
	"go-writer/service"
	"time"
)

func main() {
	//
	// file writer
	writer := gowriter.NewFileWriter(
		gowriter.WithDirectory("./testing"),
		gowriter.WithFileName("dummy_"),
		gowriter.WithFileMaxMegaByteSize(1),
		gowriter.WithFlushTime(time.Second))

	writer.Open()
	fmt.Printf("send...")
	for i := 1; i < 1000000; i++ {
		writer.Write([]byte(fmt.Sprintf("ola %d\n", i)))
	}
	fmt.Printf("sent!")
}
