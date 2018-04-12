package gomapper

import (
	"os"

	"github.com/joaosoft/go-log/service"
	"github.com/joaosoft/go-writer/service"
)

var global = make(map[string]interface{})
var log = golog.NewLog(
	golog.WithLevel(golog.InfoLevel),
	golog.WithFormatHandler(gowriter.JsonFormatHandler),
	golog.WithWriter(os.Stdout)).WithPrefixes(map[string]interface{}{
	"level":   golog.LEVEL,
	"time":    golog.TIME,
	"service": "go-mapper"})
