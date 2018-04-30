package gomapper

import (
	golog "github.com/joaosoft/go-log/app"
)

var global = make(map[string]interface{})
var log = golog.NewLogDefault("go-mapper", golog.InfoLevel)

func init() {
	global[path_key] = defaultPath
}
