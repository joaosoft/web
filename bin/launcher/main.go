package main

import (
	"fmt"
	"go-log/services"
	"os"
)

func main() {
	//
	// log to text
	fmt.Println(":: LOG TEXT")
	log := golog.NewLog(golog.WithLevel(golog.InfoLevel), golog.WithFormatHandler(golog.TextFormatHandler), golog.WithWriter(os.Stdout))
	log.With(map[string]interface{}{"level": golog.LEVEL, "time": golog.TIME}, map[string]interface{}{"service": "log"}, map[string]interface{}{"name": "joão"})

	// logging...
	log.Error("isto é uma mensagem de error")
	log.Info("isto é uma mensagem de info")
	log.Debug("isto é uma mensagem de debug")

	fmt.Println("--------------")

	//
	// log to json
	fmt.Println(":: LOG JSON")
	log = golog.NewLog(golog.WithLevel(golog.InfoLevel), golog.WithFormatHandler(golog.JsonFormatHandler), golog.WithWriter(os.Stdout))
	log.With(map[string]interface{}{"level": golog.LEVEL, "time": golog.TIME}, map[string]interface{}{"service": "log"}, map[string]interface{}{"name": "joão"})

	// logging...
	log.Errorf("isto é uma mensagem de error %s", "hello")
	log.Infof("isto é uma  mensagem de info %s ", "hi")
	log.Debugf("isto é uma mensagem de debug %s", "ehh")
}
